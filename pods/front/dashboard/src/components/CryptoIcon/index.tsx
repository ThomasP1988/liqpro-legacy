import React from 'react';
import CryptoData from 'cryptocurrency-icons/manifest.json';


export function CryptoIcon({ symbol, className }: { symbol: string, className?: string }) {

    const crypto = CryptoData.find(i => i.symbol.toLowerCase() === symbol.toLowerCase())
    
    return  <img className={className} alt={symbol || "crypto"} src={`/images/crypto/cryptocurrency-icons/svg/color/${crypto?.symbol?.toLowerCase()}.svg`} />;

}
