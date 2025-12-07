import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { Cloud, Shield, Zap, Globe, Server, LogOut } from 'lucide-react';
import { useTranslation } from 'react-i18next';

// GitHub Icon SVG
const GitHubIcon = ({ className = "w-5 h-5" }: { className?: string }) => (
  <svg className={className} viewBox="0 0 24 24" fill="currentColor">
    <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/>
  </svg>
);

// Google Icon SVG
const GoogleIcon = ({ className = "w-5 h-5" }: { className?: string }) => (
  <svg className={className} viewBox="0 0 24 24">
    <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
    <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
    <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
    <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
  </svg>
);

// User interface
interface CloudUser {
  id: string;
  username: string;
  email?: string;
  avatar?: string;
  provider: 'github' | 'google';
  plan: 'free' | 'pro';
  serverCount: number;
  serverLimit: number;
}

// Check if running in development mode (localhost)
const isDev = typeof window !== 'undefined' && (
  window.location.hostname === 'localhost' || 
  window.location.hostname === '127.0.0.1'
);

export default function CloudPage() {
  const { t } = useTranslation();
  const [user, setUser] = useState<CloudUser | null>(null);
  const [oauthLoading, setOauthLoading] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Check for OAuth callback and stored session
  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const provider = params.get('provider');
    const username = params.get('user');

    if (provider && username) {
      // OAuth callback - store user info
      const newUser: CloudUser = {
        id: btoa(username),
        username: username,
        email: provider === 'google' ? username : undefined,
        provider: provider as 'github' | 'google',
        plan: 'free',
        serverCount: 0,
        serverLimit: 5,
      };
      
      localStorage.setItem('vstats_cloud_user', JSON.stringify(newUser));
      setUser(newUser);
      
      // Clean URL
      window.history.replaceState({}, '', '/cloud');
    } else {
      // Check for existing session
      const storedUser = localStorage.getItem('vstats_cloud_user');
      if (storedUser) {
        try {
          setUser(JSON.parse(storedUser));
        } catch (e) {
          localStorage.removeItem('vstats_cloud_user');
        }
      }
    }
    
    setIsLoading(false);
  }, []);

  const handleOAuthLogin = (provider: 'github' | 'google') => {
    setOauthLoading(provider);
    
    // In development mode, simulate OAuth login
    if (isDev) {
      setTimeout(() => {
        const mockUsername = provider === 'github' ? 'dev-user' : 'dev@example.com';
        const newUser: CloudUser = {
          id: btoa(mockUsername),
          username: mockUsername,
          email: provider === 'google' ? mockUsername : undefined,
          provider: provider,
          plan: 'free',
          serverCount: 0,
          serverLimit: 5,
        };
        
        localStorage.setItem('vstats_cloud_user', JSON.stringify(newUser));
        setUser(newUser);
        setOauthLoading(null);
      }, 1000);
      return;
    }
    
    // Production: use OAuth proxy
    const state = btoa(JSON.stringify({
      redirect_uri: window.location.origin + '/cloud',
      timestamp: Date.now()
    })).replace(/[+/=]/g, (m) => ({ '+': '-', '/': '_', '=': '' }[m] || m));
    
    const oauthProxyUrl = 'https://auth.vstats.zsoft.cc';
    const redirectUri = encodeURIComponent(window.location.origin + '/cloud');
    const oauthUrl = `${oauthProxyUrl}/oauth/${provider}?redirect_uri=${redirectUri}&state=${state}`;
    window.location.href = oauthUrl;
  };

  const handleLogout = () => {
    localStorage.removeItem('vstats_cloud_user');
    setUser(null);
  };

  if (isLoading) {
    return (
      <div className="pt-20 min-h-screen flex items-center justify-center">
        <div className="w-8 h-8 border-2 border-violet-500 border-t-transparent rounded-full animate-spin" />
      </div>
    );
  }

  // Logged in - Show Coming Soon Dashboard
  if (user) {
    return (
      <div className="pt-20 min-h-screen bg-slate-50 dark:bg-slate-900">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
          {/* User Info Card */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="glass-card rounded-2xl p-8 mb-8"
          >
            <div className="flex flex-col sm:flex-row justify-between items-center gap-6">
              <div className="flex items-center gap-4">
                <div className="w-16 h-16 rounded-xl bg-gradient-to-br from-violet-500 to-purple-600 flex items-center justify-center text-white text-2xl font-bold shadow-lg shadow-violet-500/30">
                  {user.username.charAt(0).toUpperCase()}
                </div>
                <div className="text-center sm:text-left">
                  <h1 className="text-2xl font-bold dark:text-white flex items-center gap-2 justify-center sm:justify-start">
                    {user.username}
                    {user.provider === 'github' ? (
                      <GitHubIcon className="w-5 h-5 text-slate-400" />
                    ) : (
                      <GoogleIcon className="w-5 h-5" />
                    )}
                  </h1>
                  <p className="text-slate-500 text-sm">
                    {t('cloud.loginWith', { provider: user.provider === 'github' ? 'GitHub' : 'Google' })}
                  </p>
                </div>
              </div>
              <button
                onClick={handleLogout}
                className="btn btn-outline flex items-center gap-2"
              >
                <LogOut className="w-4 h-4" />
                {t('cloud.logout')}
              </button>
            </div>
          </motion.div>

          {/* Coming Soon Card */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 }}
            className="glass-card rounded-2xl p-12 text-center"
          >
            {/* Animated Icon */}
            <div className="relative w-24 h-24 mx-auto mb-8">
              <div className="absolute inset-0 rounded-2xl bg-violet-500/20 animate-ping" />
              <div className="relative w-24 h-24 rounded-2xl bg-gradient-to-br from-violet-500 to-purple-600 flex items-center justify-center shadow-lg shadow-violet-500/30">
                <Cloud className="w-12 h-12 text-white" />
              </div>
            </div>

            {/* Coming Soon Badge */}
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-violet-100 dark:bg-violet-900/30 text-violet-600 dark:text-violet-400 text-sm font-medium mb-6">
              <span className="relative flex h-2 w-2">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-violet-400 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-2 w-2 bg-violet-500"></span>
              </span>
              {t('cloud.inDevelopment')}
            </div>

            <h2 className="text-3xl font-bold mb-4 dark:text-white">
              {t('cloud.comingSoonTitle')}
            </h2>
            
            <p className="text-slate-500 text-lg mb-8 max-w-md mx-auto">
              {t('cloud.comingSoonDesc')}
            </p>

            {/* Features Preview */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-8">
              {[
                { icon: Server, labelKey: 'multiServer' },
                { icon: Shield, labelKey: 'enterpriseSecurityShort' },
                { icon: Zap, labelKey: 'realtime' },
              ].map((feature) => (
                <div key={feature.labelKey} className="p-4 rounded-xl bg-slate-100 dark:bg-slate-800/50">
                  <feature.icon className="w-6 h-6 text-violet-500 mx-auto mb-2" />
                  <div className="text-sm font-medium dark:text-white">{t(`cloud.features.${feature.labelKey}`)}</div>
                </div>
              ))}
            </div>

            {/* Links */}
            <div className="mt-8 pt-8 border-t border-slate-200 dark:border-slate-700">
              <p className="text-sm text-slate-500 mb-4">
                {t('cloud.waitingMessage')}
              </p>
              <div className="flex flex-wrap gap-4 justify-center">
                <a href="/" className="text-violet-600 dark:text-violet-400 hover:underline text-sm font-medium">
                  {t('cloud.learnMore')}
                </a>
                <a href="/cli" className="text-violet-600 dark:text-violet-400 hover:underline text-sm font-medium">
                  {t('cloud.cliGuide')}
                </a>
                <a href="https://github.com/zsai001/vstats" target="_blank" rel="noreferrer" className="text-violet-600 dark:text-violet-400 hover:underline text-sm font-medium">
                  {t('cloud.githubOpenSource')}
                </a>
              </div>
            </div>
          </motion.div>
        </div>
      </div>
    );
  }

  // Not logged in - Show Login Page
  return (
    <div className="pt-20 min-h-screen">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
        <div className="text-center mb-16">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="inline-flex items-center justify-center w-20 h-20 rounded-2xl bg-violet-100 dark:bg-violet-900/30 text-violet-500 mb-8"
          >
            <Cloud className="w-10 h-10" />
          </motion.div>
          <motion.h1
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 }}
            className="text-5xl font-bold mb-6 dark:text-white"
          >
            {t('cloud.title')}
          </motion.h1>
          <motion.p
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
            className="text-xl text-slate-500 max-w-2xl mx-auto mb-8"
          >
            {t('cloud.subtitle')}
          </motion.p>
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.3 }}
            className="mt-8 flex flex-col gap-4 justify-center items-center"
          >
            {/* Coming Soon Badge */}
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-violet-100 dark:bg-violet-900/30 text-violet-600 dark:text-violet-400 text-sm font-medium">
              <span className="relative flex h-2 w-2">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-violet-400 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-2 w-2 bg-violet-500"></span>
              </span>
              {t('cloud.comingSoon')}
            </div>
            
            <div className="flex flex-col sm:flex-row gap-4 justify-center items-center">
              <button
                onClick={() => handleOAuthLogin('github')}
                disabled={oauthLoading !== null}
                className="group w-full sm:w-auto py-3 px-6 rounded-xl bg-slate-900 hover:bg-slate-800 disabled:bg-slate-700 border border-slate-800 hover:border-slate-700 disabled:border-slate-700 font-medium transition-all duration-200 flex items-center justify-center gap-3 shadow-lg shadow-slate-900/30 hover:-translate-y-0.5 disabled:cursor-not-allowed text-white"
              >
                {oauthLoading === 'github' ? (
                  <>
                    <div className="w-5 h-5 border-2 border-white/20 border-t-white rounded-full animate-spin" />
                    <span>{t('cloud.loggingIn')}</span>
                  </>
                ) : (
                  <>
                    <GitHubIcon className="w-5 h-5" />
                    <span>{t('cloud.loginWithGitHub')}</span>
                  </>
                )}
              </button>
              <button
                onClick={() => handleOAuthLogin('google')}
                disabled={oauthLoading !== null}
                className="group w-full sm:w-auto py-3 px-6 rounded-xl bg-white hover:bg-slate-50 disabled:bg-slate-100 border border-slate-200 hover:border-slate-300 disabled:border-slate-200 text-slate-700 font-medium transition-all duration-200 flex items-center justify-center gap-3 shadow-lg shadow-emerald-500/10 hover:-translate-y-0.5 disabled:cursor-not-allowed"
              >
                {oauthLoading === 'google' ? (
                  <>
                    <div className="w-5 h-5 border-2 border-slate-200 border-t-slate-500 rounded-full animate-spin" />
                    <span>{t('cloud.loggingIn')}</span>
                  </>
                ) : (
                  <>
                    <GoogleIcon className="w-5 h-5" />
                    <span>{t('cloud.loginWithGoogle')}</span>
                  </>
                )}
              </button>
            </div>
            
            <p className="text-sm text-slate-400 mt-2">
              {t('cloud.devMessage')}
            </p>
          </motion.div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center mb-20">
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: 0.4 }}
            className="space-y-8"
          >
            {[
              { icon: Shield, titleKey: 'enterpriseSecurity', descKey: 'enterpriseSecurity' },
              { icon: Globe, titleKey: 'globalEdge', descKey: 'globalEdge' },
              { icon: Zap, titleKey: 'dataRetention', descKey: 'dataRetention' },
            ].map((feature) => (
              <div key={feature.titleKey} className="flex gap-4">
                <div className="w-12 h-12 rounded-xl bg-violet-100 dark:bg-violet-900/30 flex items-center justify-center text-violet-500 shrink-0">
                  <feature.icon className="w-6 h-6" />
                </div>
                <div>
                  <h3 className="text-xl font-bold mb-2 dark:text-white">{t(`cloud.features.${feature.titleKey}.title`)}</h3>
                  <p className="text-slate-500">{t(`cloud.features.${feature.descKey}.desc`)}</p>
                </div>
              </div>
            ))}
          </motion.div>
          
          <motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: 0.5 }}
            className="glass-card p-8 rounded-2xl border-t-4 border-violet-500"
          >
            <h3 className="text-2xl font-bold mb-6 dark:text-white">{t('cloud.pricing.title')}</h3>
            <div className="space-y-4">
              <div className="p-4 rounded-xl bg-slate-50 dark:bg-slate-800/50 border border-slate-200 dark:border-slate-700 flex justify-between items-center">
                <div>
                  <div className="font-bold dark:text-white">{t('cloud.pricing.free.name')}</div>
                  <div className="text-sm text-slate-500">{t('cloud.pricing.free.desc')}</div>
                </div>
                <div className="font-bold text-xl dark:text-white">{t('cloud.pricing.free.price')}</div>
              </div>
              <div className="p-4 rounded-xl bg-violet-50 dark:bg-violet-900/10 border border-violet-200 dark:border-violet-800 flex justify-between items-center">
                <div>
                  <div className="font-bold text-violet-700 dark:text-violet-300">{t('cloud.pricing.pro.name')}</div>
                  <div className="text-sm text-violet-600 dark:text-violet-400">{t('cloud.pricing.pro.desc')}</div>
                </div>
                <div className="font-bold text-xl text-violet-700 dark:text-violet-300">{t('cloud.pricing.pro.price')}<span className="text-sm font-normal">{t('cloud.pricing.pro.perServer')}</span></div>
              </div>
            </div>
          </motion.div>
        </div>
      </div>
    </div>
  );
}
