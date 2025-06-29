import './App.css'
export const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/v1"
import './pages/Login.tsx'
import Login from './pages/Login.tsx'
function App() {

  return (
    <>
      <div>App Home Screen</div>
      <Login/>
    </>
  )
}

export default App
