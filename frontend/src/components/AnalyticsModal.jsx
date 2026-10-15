import React, { useEffect, useState } from 'react';
import { X, BarChart3, MousePointerClick, Calendar, Globe, Monitor, Compass, Loader2 } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { getUrlAnalytics } from '../services/api';

export default function AnalyticsModal({ shortCode, onClose }) {
  const { t, i18n } = useTranslation();
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!shortCode) return;
    setLoading(true);
    setError('');

    getUrlAnalytics(shortCode)
      .then((res) => {
        setData(res);
      })
      .catch((err) => {
        setError(err.message || t('analyticsModal.loadError'));
      })
      .finally(() => {
        setLoading(false);
      });
  }, [shortCode, t]);

  if (!shortCode) return null;

  // Helper to render proportional metric bars
  const renderBars = (obj) => {
    if (!obj || Object.keys(obj).length === 0) {
      return <p className="text-xs text-slate-500 italic">{t('analyticsModal.noData')}</p>;
    }

    const total = Object.values(obj).reduce((acc, curr) => acc + curr, 0) || 1;

    return (
      <div className="space-y-2 mt-2">
        {Object.entries(obj).map(([key, count]) => {
          const percent = Math.round((count / total) * 100);
          return (
            <div key={key} className="space-y-1">
              <div className="flex justify-between text-xs">
                <span className="text-slate-300 font-medium truncate max-w-[150px]">{key}</span>
                <span className="text-slate-400 font-mono">
                  {count} ({percent}%)
                </span>
              </div>
              <div className="w-full bg-slate-900 rounded-full h-1.5 overflow-hidden">
                <div
                  className="bg-indigo-500 h-1.5 rounded-full transition-all duration-500"
                  style={{ width: `${percent}%` }}
                />
              </div>
            </div>
          );
        })}
      </div>
    );
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm animate-fadeIn">
      <div className="bg-slate-800 border border-slate-700 w-full max-w-2xl rounded-2xl shadow-2xl overflow-hidden flex flex-col max-h-[90vh]">
        {/* Header */}
        <div className="p-5 border-b border-slate-700/80 flex items-center justify-between bg-slate-900/50">
          <div className="flex items-center gap-2.5">
            <div className="p-2 bg-purple-500/10 text-purple-400 rounded-xl border border-purple-500/20">
              <BarChart3 className="w-5 h-5" />
            </div>
            <div>
              <h3 className="text-lg font-bold text-white flex items-center gap-2">
                {t('analyticsModal.title')} <span className="font-mono text-indigo-400">/r/{shortCode}</span>
              </h3>
              {data && (
                <p className="text-xs text-slate-400 truncate max-w-sm sm:max-w-md">
                  {data.original_url}
                </p>
              )}
            </div>
          </div>

          <button
            onClick={onClose}
            className="p-1.5 text-slate-400 hover:text-white bg-slate-800 hover:bg-slate-700 rounded-lg transition-colors cursor-pointer"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6 overflow-y-auto space-y-6">
          {loading ? (
            <div className="py-16 text-center text-slate-400 flex flex-col items-center gap-3">
              <Loader2 className="w-8 h-8 animate-spin text-indigo-400" />
              <p className="text-sm">{t('analyticsModal.loading')}</p>
            </div>
          ) : error ? (
            <div className="py-12 text-center text-rose-400">
              <p>{error}</p>
            </div>
          ) : data ? (
            <>
              {/* Stat Cards */}
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                <div className="bg-slate-900/70 border border-slate-700/60 p-4 rounded-xl">
                  <span className="text-xs text-slate-400 flex items-center gap-1.5 mb-1">
                    <MousePointerClick className="w-3.5 h-3.5 text-indigo-400" />
                    {t('analyticsModal.totalClicks')}
                  </span>
                  <div className="text-2xl font-black text-white font-mono">{data.total_clicks}</div>
                </div>

                <div className="bg-slate-900/70 border border-slate-700/60 p-4 rounded-xl">
                  <span className="text-xs text-slate-400 flex items-center gap-1.5 mb-1">
                    <Globe className="w-3.5 h-3.5 text-emerald-400" />
                    {t('analyticsModal.trafficSources')}
                  </span>
                  <div className="text-2xl font-black text-white font-mono">
                    {Object.keys(data.top_referrers || {}).length}
                  </div>
                </div>

                <div className="bg-slate-900/70 border border-slate-700/60 p-4 rounded-xl">
                  <span className="text-xs text-slate-400 flex items-center gap-1.5 mb-1">
                    <Calendar className="w-3.5 h-3.5 text-pink-400" />
                    {t('analyticsModal.lastClick')}
                  </span>
                  <div className="text-sm font-semibold text-slate-200 mt-1 truncate">
                    {data.last_clicked
                      ? new Date(data.last_clicked).toLocaleTimeString(i18n.language, {
                          hour: '2-digit',
                          minute: '2-digit',
                        })
                      : t('analyticsModal.noVisits')}
                  </div>
                </div>
              </div>

              {/* Detailed breakdowns */}
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                {/* Referrers */}
                <div className="bg-slate-900/60 border border-slate-700/50 p-4 rounded-xl">
                  <h4 className="text-xs uppercase tracking-wider font-bold text-slate-400 flex items-center gap-1.5 mb-2">
                    <Globe className="w-3.5 h-3.5 text-indigo-400" />
                    {t('analyticsModal.referrers')}
                  </h4>
                  {renderBars(data.top_referrers)}
                </div>

                {/* Devices */}
                <div className="bg-slate-900/60 border border-slate-700/50 p-4 rounded-xl">
                  <h4 className="text-xs uppercase tracking-wider font-bold text-slate-400 flex items-center gap-1.5 mb-2">
                    <Monitor className="w-3.5 h-3.5 text-emerald-400" />
                    {t('analyticsModal.devices')}
                  </h4>
                  {renderBars(data.top_devices)}
                </div>

                {/* Browsers */}
                <div className="bg-slate-900/60 border border-slate-700/50 p-4 rounded-xl">
                  <h4 className="text-xs uppercase tracking-wider font-bold text-slate-400 flex items-center gap-1.5 mb-2">
                    <Compass className="w-3.5 h-3.5 text-purple-400" />
                    {t('analyticsModal.browsers')}
                  </h4>
                  {renderBars(data.top_browsers)}
                </div>

                {/* Clicks over time */}
                <div className="bg-slate-900/60 border border-slate-700/50 p-4 rounded-xl">
                  <h4 className="text-xs uppercase tracking-wider font-bold text-slate-400 flex items-center gap-1.5 mb-2">
                    <Calendar className="w-3.5 h-3.5 text-pink-400" />
                    {t('analyticsModal.recentActivity')}
                  </h4>
                  {renderBars(data.clicks_by_date)}
                </div>
              </div>
            </>
          ) : null}
        </div>
      </div>
    </div>
  );
}
