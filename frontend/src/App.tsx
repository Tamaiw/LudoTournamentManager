import { Routes, Route } from 'react-router-dom'

function App() {
  return (
    <div className="min-h-screen bg-gray-100">
      <Routes>
        <Route path="/" element={<div>Welcome to Ludo Tournament Manager</div>} />
        <Route path="/tournaments" element={<div>Tournaments</div>} />
        <Route path="/leagues" element={<div>Leagues</div>} />
        <Route path="/login" element={<div>Login</div>} />
      </Routes>
    </div>
  )
}

export default App