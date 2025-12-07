import React from 'react';
import { Routes, Route } from 'react-router-dom';
import Navbar from './components/Navbar';
import Footer from './components/Footer';
import Home from './pages/Home';
import CLI from './pages/CLI';
import Cloud from './pages/Cloud';
import Documentation from './pages/Documentation';

function App() {
  return (
    <div className="min-h-screen bg-slate-50 dark:bg-slate-900 text-slate-900 dark:text-slate-50 font-sans selection:bg-sky-500/30">
      <Navbar />
      <main>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/cli" element={<CLI />} />
          <Route path="/cloud" element={<Cloud />} />
          <Route path="/docs" element={<Documentation />} />
        </Routes>
      </main>
      <Footer />
    </div>
  );
}

export default App;
