import { useState, useRef, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

interface LanguageSwitcherProps {
  isDark: boolean;
}

const languages = [
  { code: 'en', name: 'English', flag: '🇺🇸', short: 'EN' },
  { code: 'zh', name: '中文', flag: '🇨🇳', short: '中' },
  { code: 'ja', name: '日本語', flag: '🇯🇵', short: '日' },
  { code: 'ko', name: '한국어', flag: '🇰🇷', short: '한' },
  { code: 'ru', name: 'Русский', flag: '🇷🇺', short: 'RU' },
  { code: 'fr', name: 'Français', flag: '🇫🇷', short: 'FR' },
  { code: 'de', name: 'Deutsch', flag: '🇩🇪', short: 'DE' },
  { code: 'es', name: 'Español', flag: '🇪🇸', short: 'ES' },
  { code: 'pt', name: 'Português', flag: '🇧🇷', short: 'PT' },
];

export function LanguageSwitcher({ isDark }: LanguageSwitcherProps) {
  const { i18n } = useTranslation();
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  const currentLang = languages.find(l => l.code === i18n.language || i18n.language.startsWith(l.code)) || languages[0];

  const handleLanguageChange = (code: string) => {
    i18n.changeLanguage(code);
    setIsOpen(false);
  };

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  return (
    <div className="relative" ref={dropdownRef}>
      <button
        onClick={() => setIsOpen(!isOpen)}
        className={`vps-btn ${isDark ? 'vps-btn-outline-dark' : 'vps-btn-outline-light'} p-2.5 flex items-center gap-1.5`}
        title={`Language: ${currentLang.name}`}
      >
        <span className="text-sm">{currentLang.flag}</span>
        <span className="text-xs font-medium">{currentLang.short}</span>
        <svg className={`w-3 h-3 transition-transform ${isOpen ? 'rotate-180' : ''}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
        </svg>
      </button>

      {isOpen && (
        <div className={`absolute right-0 mt-2 py-2 w-40 rounded-xl shadow-xl z-50 border ${
          isDark 
            ? 'bg-gray-900/95 border-white/10 backdrop-blur-xl' 
            : 'bg-white/95 border-gray-200 backdrop-blur-xl'
        }`}>
          {languages.map((lang) => (
            <button
              key={lang.code}
              onClick={() => handleLanguageChange(lang.code)}
              className={`w-full px-4 py-2 text-left flex items-center gap-3 transition-colors ${
                currentLang.code === lang.code
                  ? isDark
                    ? 'bg-emerald-500/20 text-emerald-400'
                    : 'bg-emerald-500/10 text-emerald-600'
                  : isDark
                    ? 'text-gray-300 hover:bg-white/5'
                    : 'text-gray-700 hover:bg-gray-100'
              }`}
            >
              <span className="text-base">{lang.flag}</span>
              <span className="text-sm">{lang.name}</span>
              {currentLang.code === lang.code && (
                <svg className="w-4 h-4 ml-auto text-emerald-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
              )}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}

