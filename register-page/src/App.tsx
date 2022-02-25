import './App.css';

import React from 'react';
import { Route, Routes, Navigate } from 'react-router-dom';

import Header from './components/Header';
import Body from './components/Body';
import Description from './components/Description';
import Footer from './components/Footer';
import Redirect from './components/Redirect';
import Commands from './components/Commands';

function App() {
  return (
    <div className="App">

      <Header />

      <Routes>

        <Route path='/' element={<><Description /> <Body /></>} />
        <Route path='/redirect' element={<Redirect />} />
        <Route path='/commands' element={<Commands />} />
        <Route path='*' element={<Navigate to='/'/>}/>

      </Routes>

      <Footer />
    </div>
  );
}

export default App;
