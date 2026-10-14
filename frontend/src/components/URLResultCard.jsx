import React, { useState } from 'react';
import { Copy, Check, QrCode, BarChart3, ExternalLink, Download } from 'lucide-react';
import { getQrCodeUrl } from '../services/api';

export default function URLResultCard({ result, onOpenAnalytics }) {
  const [copied, setCopied] = useState(false);

  if (!result) return null;

  const handleCopy = () => {
    navigator.clipboard.writeText(result.short_url);
    setCopied(true);
    setTimeout(() => setCopied(false), 2500);
  };

  const qrImageUrl = getQrCodeUrl(result.short_code, 300);

  const handleDownloadQR = async () => {
    try {
      const response = await fetch(qrImageUrl);
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `qr-${result.short_code}.png`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    } catch (err) {
      console.error('Error downloading QR code:', err);
    }
  };

  return (
    <div className="w-full max-w-3xl mx-auto mt-6 bg-slate-800/80 border border-indigo-500/30 rounded-2xl p-6 shadow-2xl shadow-indigo-950/40 animate-fadeIn">
      <div className="flex flex-col md:flex-row items-center gap-6 justify-between">
        {/* Left: URL Details & Actions */}
        <div className="flex-1 w-full space-y-4">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 text-xs font-semibold">
            <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
            ¡Enlace acortado con éxito!
          </div>

          <div>
            <span className="text-xs text-slate-400 block mb-1">Tu enlace corto:</span>
            <div className="flex items-center gap-2 bg-slate-900/90 border border-slate-700 rounded-xl p-2.5">
              <a
                href={result.short_url}
                target="_blank"
                rel="noopener noreferrer"
                className="text-indigo-400 hover:text-indigo-300 font-mono font-medium text-base truncate flex-1 flex items-center gap-1.5"
              >
                {result.short_url}
                <ExternalLink className="w-3.5 h-3.5 shrink-0 opacity-70" />
              </a>

              <button
                onClick={handleCopy}
                className="px-3 py-1.5 bg-indigo-600 hover:bg-indigo-500 text-white rounded-lg text-xs font-medium flex items-center gap-1.5 transition-all cursor-pointer"
              >
                {copied ? (
                  <>
                    <Check className="w-3.5 h-3.5 text-emerald-300" />
                    <span>¡Copiado!</span>
                  </>
                ) : (
                  <>
                    <Copy className="w-3.5 h-3.5" />
                    <span>Copiar</span>
                  </>
                )}
              </button>
            </div>
          </div>

          <div>
            <span className="text-xs text-slate-400 block mb-1">Destino original:</span>
            <p className="text-xs text-slate-300 truncate max-w-md font-mono bg-slate-900/40 p-2 rounded-lg border border-slate-800">
              {result.original_url}
            </p>
          </div>

          <div className="pt-2 flex flex-wrap gap-3">
            <button
              onClick={() => onOpenAnalytics(result.short_code)}
              className="px-4 py-2 bg-slate-700/60 hover:bg-slate-700 text-slate-200 rounded-xl text-xs sm:text-sm font-medium flex items-center gap-2 border border-slate-600/60 transition-colors cursor-pointer"
            >
              <BarChart3 className="w-4 h-4 text-purple-400" />
              Ver Analíticas de Clics
            </button>
          </div>
        </div>

        {/* Right: QR Code Visual */}
        <div className="flex flex-col items-center p-4 bg-slate-900/90 border border-slate-700/80 rounded-xl shrink-0">
          <div className="bg-white p-2.5 rounded-lg shadow-inner">
            <img
              src={qrImageUrl}
              alt={`QR para ${result.short_code}`}
              className="w-32 h-32 md:w-36 md:h-36 object-contain"
            />
          </div>

          <button
            onClick={handleDownloadQR}
            className="mt-3 w-full py-1.5 px-3 bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-medium rounded-lg flex items-center justify-center gap-1.5 transition-colors cursor-pointer border border-slate-700"
          >
            <Download className="w-3.5 h-3.5 text-indigo-400" />
            Descargar QR
          </button>
        </div>
      </div>
    </div>
  );
}
