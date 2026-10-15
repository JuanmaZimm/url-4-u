import React, { useState } from 'react';
import { Link, Sparkles, ArrowRight, Loader2, AlertCircle } from 'lucide-react';
import { useTranslation } from 'react-i18next';

export default function URLShortenerForm({ onShortened, isLoading }) {
  const { t } = useTranslation();
  const [originalUrl, setOriginalUrl] = useState('');
  const [customAlias, setCustomAlias] = useState('');
  const [showAliasInput, setShowAliasInput] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    if (!originalUrl.trim()) {
      setError(t('form.errors.empty'));
      return;
    }

    try {
      new URL(originalUrl.startsWith('http') ? originalUrl : `https://${originalUrl}`);
    } catch {
      setError(t('form.errors.invalid'));
      return;
    }

    const formattedUrl = originalUrl.startsWith('http') ? originalUrl : `https://${originalUrl}`;

    try {
      await onShortened(formattedUrl, customAlias);
      setOriginalUrl('');
      setCustomAlias('');
      setShowAliasInput(false);
    } catch (err) {
      setError(err.message || t('form.errors.default'));
    }
  };

  return (
    <div className="w-full max-w-3xl mx-auto">
      <div className="bg-slate-800/60 border border-slate-700/80 backdrop-blur-md rounded-2xl p-6 md:p-8 shadow-xl shadow-black/20">
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-2">
              {t('form.label')}
            </label>
            <div className="relative flex items-center">
              <div className="absolute left-4 text-slate-400">
                <Link className="w-5 h-5" />
              </div>
              <input
                type="text"
                value={originalUrl}
                onChange={(e) => {
                  setOriginalUrl(e.target.value);
                  if (error) setError('');
                }}
                placeholder={t('form.placeholder')}
                className="w-full pl-12 pr-4 py-3.5 bg-slate-900/90 border border-slate-700 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 rounded-xl text-white placeholder-slate-500 text-base transition-all outline-none"
                disabled={isLoading}
              />
            </div>
          </div>

          <div className="flex items-center justify-between pt-1">
            <button
              type="button"
              onClick={() => setShowAliasInput(!showAliasInput)}
              className="text-xs sm:text-sm text-indigo-400 hover:text-indigo-300 font-medium flex items-center gap-1.5 transition-colors"
            >
              <Sparkles className="w-4 h-4" />
              {showAliasInput ? t('form.hideAlias') : t('form.showAlias')}
            </button>
          </div>

          {showAliasInput && (
            <div className="animate-fadeIn transition-all">
              <label className="block text-xs font-medium text-slate-400 mb-1.5">
                {t('form.aliasLabel')}
              </label>
              <div className="flex rounded-xl overflow-hidden border border-slate-700 bg-slate-900/90 focus-within:border-indigo-500 focus-within:ring-2 focus-within:ring-indigo-500/20">
                <span className="px-3.5 py-2.5 bg-slate-800/80 text-slate-400 text-sm border-r border-slate-700 select-none flex items-center">
                  url4u.dev/r/
                </span>
                <input
                  type="text"
                  value={customAlias}
                  onChange={(e) => setCustomAlias(e.target.value.replace(/[^a-zA-Z0-9_-]/g, ''))}
                  placeholder={t('form.aliasPlaceholder')}
                  maxLength={30}
                  className="w-full px-3 py-2.5 bg-transparent text-white placeholder-slate-500 text-sm outline-none"
                  disabled={isLoading}
                />
              </div>
            </div>
          )}

          {error && (
            <div className="flex items-center gap-2 p-3 bg-rose-500/10 border border-rose-500/20 rounded-xl text-rose-400 text-sm">
              <AlertCircle className="w-4 h-4 shrink-0" />
              <span>{error}</span>
            </div>
          )}

          <button
            type="submit"
            disabled={isLoading}
            className="w-full py-3.5 px-6 bg-gradient-to-r from-indigo-500 via-indigo-600 to-purple-600 hover:from-indigo-600 hover:to-purple-700 active:scale-[0.99] text-white font-semibold rounded-xl shadow-lg shadow-indigo-600/30 flex items-center justify-center gap-2 transition-all disabled:opacity-60 disabled:cursor-not-allowed cursor-pointer"
          >
            {isLoading ? (
              <>
                <Loader2 className="w-5 h-5 animate-spin" />
                <span>{t('form.loadingBtn')}</span>
              </>
            ) : (
              <>
                <span>{t('form.shortenBtn')}</span>
                <ArrowRight className="w-5 h-5" />
              </>
            )}
          </button>
        </form>
      </div>
    </div>
  );
}
