import React from 'react';
import { X, Download } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { getQrCodeUrl } from '../services/api';

export default function QRModal({ link, onClose }) {
  const { t } = useTranslation();
  if (!link) return null;

  const qrImageUrl = getQrCodeUrl(link.short_code, 400);

  const handleDownload = async () => {
    try {
      const response = await fetch(qrImageUrl);
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `qr-${link.short_code}.png`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    } catch (err) {
      console.error('Download error:', err);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm animate-fadeIn">
      <div className="bg-slate-800 border border-slate-700 w-full max-w-sm rounded-2xl shadow-2xl p-6 flex flex-col items-center text-center">
        <div className="w-full flex justify-end">
          <button
            onClick={onClose}
            className="p-1.5 text-slate-400 hover:text-white bg-slate-900/60 hover:bg-slate-700 rounded-lg transition-colors cursor-pointer"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        <div className="p-3 bg-white rounded-2xl shadow-xl mt-2 mb-4">
          <img
            src={qrImageUrl}
            alt={`QR ${link.short_code}`}
            className="w-56 h-56 object-contain"
          />
        </div>

        <h3 className="text-base font-bold text-white font-mono">{link.short_url}</h3>
        <p className="text-xs text-slate-400 truncate max-w-xs mt-1 mb-5">
          {link.original_url}
        </p>

        <button
          onClick={handleDownload}
          className="w-full py-2.5 px-4 bg-indigo-600 hover:bg-indigo-500 text-white font-medium rounded-xl text-sm flex items-center justify-center gap-2 shadow-lg shadow-indigo-600/30 transition-colors cursor-pointer"
        >
          <Download className="w-4 h-4" />
          {t('qrModal.downloadBtn')}
        </button>
      </div>
    </div>
  );
}
