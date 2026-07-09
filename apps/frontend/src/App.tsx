import { Route, Routes } from 'react-router'
import './App.css'
import AuthNavigator from './components/navigators/AuthNavigator'
import Login from './components/auth/Login'
import Create from './components/auth/Create'

function App() {
  console.log("Hello!")
  return (
    <>
      <Routes>
        <Route element={<AuthNavigator />}>
          <Route path='/login' element={<Login />} />
          <Route path='/create' element={<Create />} />
        </Route>
      </Routes>
    </>
  )
}

export default App
