import React from 'react';
import { NavLink } from 'react-router';

function App() {
  return (
    <div className="App">
      <NavLink to={"/signin"}> Signin </NavLink>
      <NavLink to={"/lender"}> lender </NavLink>
    </div>
  );
}

export default App;
