import { createContext, useContext, useState, useEffect, type ReactNode } from 'react';

interface OAuthProviders {
  github?: boolean;
  google?: boolean;
}

interface AuthContextType {
  isAuthenticated: boolean;
  token: string | null;
  login: (password: string) => Promise<boolean>;
  logout: () => void;
  isLoading: boolean;
  oauthProviders: OAuthProviders;
  startOAuthLogin: (provider: 'github' | 'google') => Promise<void>;
  handleOAuthCallback: (token: string, expiresAt: number, provider: string, user: string) => void;
  oauthUser: string | null;
  oauthProvider: string | null;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(() => localStorage.getItem('vstats_token'));
  const [isLoading, setIsLoading] = useState(true);
  const [oauthProviders, setOauthProviders] = useState<OAuthProviders>({});
  const [oauthUser, setOauthUser] = useState<string | null>(() => localStorage.getItem('vstats_oauth_user'));
  const [oauthProvider, setOauthProvider] = useState<string | null>(() => localStorage.getItem('vstats_oauth_provider'));

  // Fetch available OAuth providers
  useEffect(() => {
    const fetchProviders = async () => {
      try {
        const res = await fetch('/api/auth/oauth/providers');
        if (res.ok) {
          const data = await res.json();
          setOauthProviders(data.providers || {});
        }
      } catch (e) {
        console.error('Failed to fetch OAuth providers', e);
      }
    };
    fetchProviders();
  }, []);

  useEffect(() => {
    // Verify token on mount
    const verifyToken = async () => {
      if (token) {
        try {
          const res = await fetch('/api/auth/verify', {
            headers: { Authorization: `Bearer ${token}` }
          });
          if (!res.ok) {
            setToken(null);
            setOauthUser(null);
            setOauthProvider(null);
            localStorage.removeItem('vstats_token');
            localStorage.removeItem('vstats_oauth_user');
            localStorage.removeItem('vstats_oauth_provider');
          }
        } catch {
          // Keep token if server is unreachable
        }
      }
      setIsLoading(false);
    };
    verifyToken();
  }, [token]);

  const login = async (password: string): Promise<boolean> => {
    try {
      const res = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ password })
      });
      
      if (res.ok) {
        const data = await res.json();
        setToken(data.token);
        setOauthUser(null);
        setOauthProvider(null);
        localStorage.setItem('vstats_token', data.token);
        localStorage.removeItem('vstats_oauth_user');
        localStorage.removeItem('vstats_oauth_provider');
        return true;
      }
      return false;
    } catch {
      return false;
    }
  };

  const startOAuthLogin = async (provider: 'github' | 'google'): Promise<void> => {
    try {
      const res = await fetch(`/api/auth/oauth/${provider}`);
      if (res.ok) {
        const data = await res.json();
        // Redirect to OAuth provider
        window.location.href = data.url;
      } else {
        const error = await res.json();
        throw new Error(error.error || 'Failed to start OAuth login');
      }
    } catch (e) {
      console.error('OAuth login error:', e);
      throw e;
    }
  };

  const handleOAuthCallback = (newToken: string, expiresAt: number, provider: string, user: string) => {
    setToken(newToken);
    setOauthUser(user);
    setOauthProvider(provider);
    localStorage.setItem('vstats_token', newToken);
    localStorage.setItem('vstats_oauth_user', user);
    localStorage.setItem('vstats_oauth_provider', provider);
    localStorage.setItem('vstats_token_expires', expiresAt.toString());
  };

  const logout = () => {
    setToken(null);
    setOauthUser(null);
    setOauthProvider(null);
    localStorage.removeItem('vstats_token');
    localStorage.removeItem('vstats_oauth_user');
    localStorage.removeItem('vstats_oauth_provider');
    localStorage.removeItem('vstats_token_expires');
  };

  return (
    <AuthContext.Provider value={{ 
      isAuthenticated: !!token, 
      token, 
      login, 
      logout,
      isLoading,
      oauthProviders,
      startOAuthLogin,
      handleOAuthCallback,
      oauthUser,
      oauthProvider
    }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
}
