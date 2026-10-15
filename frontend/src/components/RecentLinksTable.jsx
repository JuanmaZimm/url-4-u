import React, { useState } from 'react';
import { Copy, Check, BarChart3, QrCode, MousePointerClick, RefreshCw } from 'lucide-react';
import { useTranslation } from 'react-i18next';

export default function RecentLinksTable({ links, onRefresh, onOpenAnalytics, onOpenQR }) {
  const { t, i18n } = useTranslation();
  const [copiedCode, setCopiedCode] = useState(null);

  const handleCopy = (code, shortUrl) => {
    navigator.clipboard.writeText(shortUrl);
    setCopiedCode(code);
    setTimeout(() => setCopiedCode(null), 2000);
  };

  return (
    <div className="w-full max-w-5xl mx-auto mt-12">
      <div className="flex items-center justify-between mb-4">
        <div>
          <h2 className="text-xl font-bold text-slate-100 flex items-center gap-2">
            <span>{t('recent.title')}</span>
            <span className="text-xs bg-slate-800 text-slate-400 font-semibold px-2 py-0.5 rounded-full border border-slate-700">
              {links.length}
            </span>
          </h2>
          <p className="text-xs text-slate-400">
            {t('recent.subtitle')}
          </p>
        </div>

        <button
          onClick={onRefresh}
          title={t('recent.refresh')}
          className="p-2 text-slate-400 hover:text-white bg-slate-800 hover:bg-slate-700 rounded-lg border border-slate-700 transition-colors cursor-pointer"
        >
          <RefreshCw className="w-4 h-4" />
        </button>
      </div>

      {links.length === 0 ? (
        <div className="text-center py-12 bg-slate-800/30 border border-slate-800 rounded-2xl text-slate-500">
          <MousePointerClick className="w-10 h-10 mx-auto mb-2 opacity-40 text-indigo-400" />
          <p className="text-sm font-medium">{t('recent.emptyTitle')}</p>
          <p className="text-xs text-slate-500 mt-1">
            {t('recent.emptySubtitle')}
          </p>
        </div>
      ) : (
        <div className="overflow-x-auto bg-slate-800/40 border border-slate-700/60 rounded-2xl shadow-xl backdrop-blur">
          <table className="w-full text-left text-sm text-slate-300">
            <thead className="bg-slate-900/60 text-xs uppercase tracking-wider text-slate-400 border-b border-slate-700/60 font-semibold">
              <tr>
                <th className="py-3.5 px-4">{t('recent.thShort')}</th>
                <th className="py-3.5 px-4 hidden md:table-cell">{t('recent.thOriginal')}</th>
                <th className="py-3.5 px-4 text-center">{t('recent.thClicks')}</th>
                <th className="py-3.5 px-4 hidden sm:table-cell">{t('recent.thCreated')}</th>
                <th className="py-3.5 px-4 text-right">{t('recent.thActions')}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/80">
              {links.map((link) => (
                <tr key={link.short_code} className="hover:bg-slate-800/50 transition-colors">
                  <td className="py-3 px-4">
                    <div className="flex items-center gap-2">
                      <a
                        href={link.short_url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-indigo-400 hover:text-indigo-300 font-mono font-medium truncate max-w-[180px] sm:max-w-[220px]"
                      >
                        {link.short_url}
                      </a>
                      <button
                        onClick={() => handleCopy(link.short_code, link.short_url)}
                        className="text-slate-400 hover:text-white p-1 rounded hover:bg-slate-700 transition-colors cursor-pointer"
                        title={t('recent.copyTitle')}
                      >
                        {copiedCode === link.short_code ? (
                          <Check className="w-3.5 h-3.5 text-emerald-400" />
                        ) : (
                          <Copy className="w-3.5 h-3.5" />
                        )}
                      </button>
                    </div>
                  </td>

                  <td className="py-3 px-4 hidden md:table-cell">
                    <span className="truncate block max-w-xs text-xs text-slate-400" title={link.original_url}>
                      {link.original_url}
                    </span>
                  </td>

                  <td className="py-3 px-4 text-center">
                    <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-semibold bg-indigo-500/10 text-indigo-300 border border-indigo-500/20">
                      <MousePointerClick className="w-3 h-3" />
                      {link.click_count}
                    </span>
                  </td>

                  <td className="py-3 px-4 hidden sm:table-cell text-xs text-slate-400">
                    {new Date(link.created_at).toLocaleDateString(i18n.language, {
                      month: 'short',
                      day: 'numeric',
                      hour: '2-digit',
                      minute: '2-digit',
                    })}
                  </td>

                  <td className="py-3 px-4 text-right">
                    <div className="flex items-center justify-end gap-1.5">
                      <button
                        onClick={() => onOpenQR(link)}
                        className="p-1.5 text-slate-300 hover:text-white bg-slate-700/60 hover:bg-slate-700 rounded-lg transition-colors cursor-pointer border border-slate-600/40"
                        title={t('recent.qrTitle')}
                      >
                        <QrCode className="w-3.5 h-3.5" />
                      </button>

                      <button
                        onClick={() => onOpenAnalytics(link.short_code)}
                        className="p-1.5 text-purple-300 hover:text-white bg-purple-500/10 hover:bg-purple-600 rounded-lg border border-purple-500/20 transition-colors cursor-pointer"
                        title={t('recent.analyticsTitle')}
                      >
                        <BarChart3 className="w-3.5 h-3.5" />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
