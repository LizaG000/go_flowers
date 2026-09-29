import { Route, Routes } from 'react-router'
import Registration from './pages/auth/Registration.tsx'
import Main from './pages/Main.tsx'
import Login from './pages/auth/Login.tsx'

function App() {
  return (
    <>

      <Routes>
        <Route path="/" element={<Main />} />
        <Route path="/registration" element={<Registration />} />
        <Route path="/login" element={<Login />} />
      </Routes>
    </>
  )
}

export default App