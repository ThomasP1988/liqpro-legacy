package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/supertokens/supertokens-go/supertokens"

	"liqpro/config"

	entities "liqpro/shared/repositories/entities"

	repositories "liqpro/shared/repositories"
)

// Login payload
type Login struct {
	Email    string
	Password string
}

func main() {

	config.SetConfig()

	cookieSecure := false
	supertokens.Config(supertokens.ConfigMap{
		Hosts:          "http://192.168.49.2:30673",
		APIKey:         "7244a844-957f-49fc-9c53-0a5dbf3fc27c",
		CookieSecure:   &cookieSecure,
		RefreshAPIPath: "/auth/session/refresh",
		CookieDomain:   "https://localhost:3000",
	})

	app := fiber.New()

	app.Use(logger.New(), cors.New(cors.Config{

		AllowOrigins:     "https://127.0.0.1:3000, https://localhost:3000",
		AllowCredentials: true,
		// ExposeHeaders:    "front-token, id-refresh-token, anti-csrf",
	}))

	// TEMPORARY REGISTRATION AND LOGIN UNTIL SUPERTOKEN IMPLEMENTS LOGIN SYSTEM

	app.Post("/login", func(c *fiber.Ctx) error {
		payload := &Login{}
		if err := json.Unmarshal([]byte(c.Body()), payload); err != nil {
			c.Status(400).Send([]byte("Invalid data"))
			return err
		}
		u := &entities.User{}
		err := repositories.GetUserRepository().FindOneByMail(payload.Email, u)
		if err != nil {
			fmt.Println("error getting user", err)
			return err
		}

		if payload.Password != u.Password {
			return errors.New("Invalid mail or password")
		}

		responseWriter := httptest.NewRecorder()
		_, err = supertokens.CreateNewSession(responseWriter, u.ID)
		if err != nil {
			fmt.Println("error creating session", err)
			return err
		}

		for k := range responseWriter.Header() {
			// if k == "Set-Cookie" {
			// 	continue
			// }
			c.Context().Response.Header.Add(k, responseWriter.Header().Get(k))
		}

		httpCookies := responseWriter.Result().Cookies()
		for k := range httpCookies {

			var sameSite string

			switch httpCookies[k].SameSite {
			case http.SameSiteLaxMode:
				sameSite = "lax"
			case http.SameSiteStrictMode:
				sameSite = "strict"
			case http.SameSiteNoneMode:
				sameSite = "none"
			case http.SameSiteDefaultMode:
				sameSite = "default"
			}

			newCookie := &fiber.Cookie{
				Name:     httpCookies[k].Name,
				Value:    httpCookies[k].Value,
				Path:     httpCookies[k].Path,
				Domain:   httpCookies[k].Domain,
				MaxAge:   httpCookies[k].MaxAge,
				Expires:  httpCookies[k].Expires,
				HTTPOnly: false,
				SameSite: sameSite,
				Secure:   httpCookies[k].Secure,
			}

			c.Cookie(newCookie)
		}

		return c.SendString("Success")
	})

	app.Post("/register", func(c *fiber.Ctx) error {

		u := &entities.User{}

		if err := json.Unmarshal([]byte(c.Body()), u); err != nil {
			c.Status(400).Send([]byte("Invalid data"))
			return err
		}

		bv := []byte(u.Email)
		hasher := sha1.New()
		hasher.Write(bv)
		sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
		u.ID = sha

		err := repositories.GetUserRepository().Create(u)
		if err != nil {
			fmt.Println("error creating user", err)
			return err
		}
		fmt.Println("done")
		return c.Status(200).SendString("Success")
	})

	// END TEMPORARY

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World !")
	})

	app.Post("/auth/session/refresh", adaptor.HTTPHandlerFunc(supertokens.Middleware(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("success"))
	})))

	// protected := app.Group("/protected", func(c *fiber.Ctx) error {

	// 	test := c.Context().Value(0)
	// 	fmt.Println("test", test)

	// 	return c.Next()
	// })

	protected := app.Group("/protected", adaptor.HTTPMiddleware(func(h http.Handler) http.Handler {
		mw := supertokens.Middleware(func(w http.ResponseWriter, r *http.Request) {
			session := supertokens.GetSessionFromRequest(r)
			// fmt.Println("session", session)
			userID := session.GetUserID()
			r.Header.Add("userID", userID)
			h.ServeHTTP(w, r)
		})
		return mw
	}))

	// TEMPORARY
	app.Get("/logout", adaptor.HTTPMiddleware(func(h http.Handler) http.Handler {

		mw := supertokens.Middleware(func(w http.ResponseWriter, r *http.Request) {
			session := supertokens.GetSessionFromRequest(r)
			err := session.RevokeSession()
			if err != nil {
				supertokens.HandleErrorAndRespond(err, w)
				return
			}

			w.Write([]byte("Logout successful"))
		})

		return mw

	}))
	// END TEMPORARY

	// START GRAPHQL

	protected.Post("/graphql", func(c *fiber.Ctx) error {
		p := &PostData{}
		if err := json.Unmarshal([]byte(c.Body()), p); err != nil {
			c.Status(400).Send([]byte("Invalid data"))
			return err
		}

		result := executeQuery(p, c)

		jsonByteResult, err := json.Marshal(result)
		if err != nil {
			c.Status(400).Send([]byte("Error formatting JSON result"))
		}

		c.Status(200).Send(jsonByteResult)
		return nil
	})

	// END GRAPHQL

	// _ = app.Listen(":8081")
	_ = app.ListenTLS(":8081", "../../../shared/local_ssl/localhost.crt", "../../../shared/local_ssl/localhost.key")
	// _ = app.ListenTLS(":8081", "", "")
}
