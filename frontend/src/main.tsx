import React from 'react'
import {createRoot} from 'react-dom/client'
import './style.css'
import 'antd/dist/reset.css';
import App from './App'
import Home from './Home/Home';
import { MemoryRouter as Router , Routes, Route} from 'react-router-dom';

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
    <React.StrictMode>
      <Router>
        <Routes>
          <Route path="/" element={<App/>} />
          <Route path="/home" element={<Home/>} />
        </Routes>
 	    </Router>
    </React.StrictMode>
)
