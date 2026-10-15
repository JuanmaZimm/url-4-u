import React from 'react';
import { Link2, Github, Globe } from 'lucide-react';
import { useTranslation } from 'react-i18next';

export default function Navbar() {
  const { t, i18n } = useTranslation();
  const currentLang = i18n.resolvedLanguage?.startsWith('es') ? 'es' : 'en';

  const toggleLanguage = () => {
    const nextLang = currentLang === 'es' ? 'en' : 'es';
    i18n.changeLanguage(nextLang);
  };

  return (
    <header className="border-b border-slate-800 bg-slate-900/80 backdrop-blur sticky top-0 z-40">
      <div className="max-w-6xl mx-auto px-4 h-16 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="bg-indigo-600 text-white p-2 rounded-xl shadow-lg shadow-indigo-500/20">
            <Link2 className="w-5 h-5" />
          </div>
          <div>
            <span className="text-xl font-bold bg-gradient-to-r from-indigo-400 via-purple-300 to-pink-400 bg-clip-text text-transparent">
              URL-4-U
            </span>
            <span className="hidden sm:inline-block ml-2 text-xs font-semibold px-2 py-0.5 rounded-full bg-indigo-500/10 text-indigo-400 border border-indigo-500/20">
              {t('navbar.version')}
            </span>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <button
            onClick={toggleLanguage}
            title={t('navbar.language')}
            aria-label={t('navbar.language')}
            className="flex items-center gap-1.5 text-xs font-medium bg-slate-800/80 hover:bg-slate-800 text-slate-300 hover:text-white px-2.5 py-1.5 rounded-lg border border-slate-700/60 transition-colors cursor-pointer"
          >
            <Globe className="w-3.5 h-3.5 text-indigo-400" />
            <span className={currentLang === 'en' ? 'font-bold text-indigo-400' : 'text-slate-400'}>
              EN
            </span>
            <span className="text-slate-600">/</span>
            <span className={currentLang === 'es' ? 'font-bold text-indigo-400' : 'text-slate-400'}>
              ES
            </span>
          </button>

          <a
            href="https://github.com/JuanmaZimm/url-4-u"
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-2 text-sm text-slate-400 hover:text-white transition-colors bg-slate-800/80 hover:bg-slate-800 px-3 py-1.5 rounded-lg border border-slate-700/60"
          >
            <Github className="w-4 h-4" />
            <span className="hidden sm:inline">{t('navbar.github')}</span>
          </a>
        </div>
      </div>
    </header>
  );
}
