import React, { useState, useEffect } from 'react';
import Navbar from './components/Navbar';
import URLShortenerForm from './components/URLShortenerForm';
import URLResultCard from './components/URLResultCard';
import RecentLinksTable from './components/RecentLinksTable';
import AnalyticsModal from './components/AnalyticsModal';
import QRModal from './components/QRModal';
import { shortenUrl, getRecentUrls } from './services/api';
import { Sparkles, ShieldCheck, Zap, BarChart2 } from 'lucide-react';

export default function App() {
  const [recentLinks, setRecentLinks] = useState([]);
  const [latestResult, setLatestResult] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [selectedAnalyticsCode, setSelectedAnalyticsCode] = useState(null);
  const [selectedQRLink, setSelectedQRLink] = useState(null);

  const fetchRecent = async () => {
    try {
      const data = await getRecentUrls(15, 0);
      setRecentLinks(data);
    } catch (err) {
      console.error('Error fetching recent URLs:', err);
    }
  };

  useEffect(() => {
    fetchRecent();
  }, []);

  const handleShorten = async (originalUrl, customAlias) => {
    setIsLoading(true);
    try {
      const result = await shortenUrl(originalUrl, customAlias);
      setLatestResult(result);
      await fetchRecent();
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col selection:bg-indigo-500 selection:text-white">
      <Navbar />

      <main className="flex-1 max-w-6xl w-full mx-auto px-4 py-12 md:py-16">
        {/* Hero Section */}
        <div className="text-center max-w-3xl mx-auto mb-10 space-y-4">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-indigo-500/10 text-indigo-400 border border-indigo-500/20 text-xs font-semibold tracking-wide uppercase">
            <Sparkles className="w-3.5 h-3.5" />
            Acortador + QR + Métricas en Vivo
          </div>

          <h1 className="text-4xl sm:text-5xl md:text-6xl font-black tracking-tight text-white leading-tight">
            Acorta tus enlaces. <br className="hidden sm:inline" />
            <span className="bg-gradient-to-r from-indigo-400 via-purple-300 to-pink-400 bg-clip-text text-transparent">
              Mide cada interacción.
            </span>
          </h1>

          <p className="text-slate-400 text-sm sm:text-base max-w-xl mx-auto">
            Crea URLs cortas con identificadores únicos o alias personalizados, genera códigos QR instantáneos para impresión o compartir, y consulta analíticas de visitas en tiempo real.
          </p>
        </div>

        {/* Shortener Form */}
        <URLShortenerForm onShortened={handleShorten} isLoading={isLoading} />

        {/* Newly Shortened URL Result */}
        {latestResult && (
          <URLResultCard
            result={latestResult}
            onOpenAnalytics={(code) => setSelectedAnalyticsCode(code)}
          />
        )}

        {/* Features banner */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 my-14 max-w-4xl mx-auto">
          <div className="bg-slate-900/50 border border-slate-800/80 p-5 rounded-2xl flex items-start gap-3.5">
            <div className="p-2.5 bg-indigo-500/10 text-indigo-400 rounded-xl border border-indigo-500/20 shrink-0">
              <Zap className="w-5 h-5" />
            </div>
            <div>
              <h4 className="text-sm font-bold text-white">Ultra Rápido</h4>
              <p className="text-xs text-slate-400 mt-1">
                Motor HTTP de baja latencia escrito en Go con redirecciones concurrentes inmediatas.
              </p>
            </div>
          </div>

          <div className="bg-slate-900/50 border border-slate-800/80 p-5 rounded-2xl flex items-start gap-3.5">
            <div className="p-2.5 bg-emerald-500/10 text-emerald-400 rounded-xl border border-emerald-500/20 shrink-0">
              <ShieldCheck className="w-5 h-5" />
            </div>
            <div>
              <h4 className="text-sm font-bold text-white">Privacidad y Seguridad</h4>
              <p className="text-xs text-slate-400 mt-1">
                IPs anonimizadas mediante hash SHA-256 para respetar la privacidad en analíticas.
              </p>
            </div>
          </div>

          <div className="bg-slate-900/50 border border-slate-800/80 p-5 rounded-2xl flex items-start gap-3.5">
            <div className="p-2.5 bg-purple-500/10 text-purple-400 rounded-xl border border-purple-500/20 shrink-0">
              <BarChart2 className="w-5 h-5" />
            </div>
            <div>
              <h4 className="text-sm font-bold text-white">Métricas Detalladas</h4>
              <p className="text-xs text-slate-400 mt-1">
                Conteo de clics, fecha, fuentes de tráfico (referrers), navegadores y tipos de dispositivo.
              </p>
            </div>
          </div>
        </div>

        {/* Recent Links Table */}
        <RecentLinksTable
          links={recentLinks}
          onRefresh={fetchRecent}
          onOpenAnalytics={(code) => setSelectedAnalyticsCode(code)}
          onOpenQR={(link) => setSelectedQRLink(link)}
        />
      </main>

      {/* Analytics Modal */}
      {selectedAnalyticsCode && (
        <AnalyticsModal
          shortCode={selectedAnalyticsCode}
          onClose={() => setSelectedAnalyticsCode(null)}
        />
      )}

      {/* QR Code Modal */}
      {selectedQRLink && (
        <QRModal
          link={selectedQRLink}
          onClose={() => setSelectedQRLink(null)}
        />
      )}

      {/* Footer */}
      <footer className="border-t border-slate-900 py-8 text-center text-xs text-slate-500">
        <p>
          URL-4-U • Construido con <span className="text-indigo-400 font-medium">Go</span>,{' '}
          <span className="text-indigo-400 font-medium">React + Vite</span> y{' '}
          <span className="text-indigo-400 font-medium">SQLite</span>.
        </p>
      </footer>
    </div>
  );
}
