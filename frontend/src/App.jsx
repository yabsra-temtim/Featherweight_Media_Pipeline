import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { ThemeProvider } from './context/ThemeContext'
import NavBar from './components/NavBar'
import Home from './pages/Home'
import Upload from './pages/Upload'
import About from './pages/About'

export default function App() {
    return (
        <ThemeProvider>
            <BrowserRouter>
                <div className="min-h-screen">
                    <NavBar />
                    <Routes>
                        <Route path="/" element={<Home />} />
                        <Route path="/upload" element={<Upload />} />
                        <Route path="/about" element={<About />} />
                    </Routes>
                </div>
            </BrowserRouter>
        </ThemeProvider>
    )
}
