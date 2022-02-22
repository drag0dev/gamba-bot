import './App.css';

import React from 'react';
import { Route, Routes } from 'react-router-dom';

import Header from './components/Header';
import Body from './components/Body';
import Description from './components/Desciprion';
import Footer from './components/Footer';

function App() {
  return (
    <div className="App">

      <Header />

      <Routes>

        <Route path='/' element={<><Description /> <Body /></>} />

      </Routes>

      <Footer />
    </div>
  );
}

export default App;
