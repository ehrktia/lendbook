import './index.css';
import OAuthSignInPage from './signin';
import reportWebVitals from './reportWebVitals';
import { BrowserRouter, Route, Routes } from 'react-router';
import App from './App';
import Lender from './lender';
import Logout from './logout';
import { createRoot } from 'react-dom/client';
import { ApolloProvider, ApolloClient, InMemoryCache } from '@apollo/client';

const apolloClient = new ApolloClient({
  uri:"http://localhost:8080/query",
  cache: new InMemoryCache(),
  })
const root = createRoot(
  document.getElementById('root') as HTMLElement
);

root.render(
  <ApolloProvider client={apolloClient}>
  <BrowserRouter>
    <Routes>
      <Route path='/' element={<App />} />
      <Route path="/signin" element={<OAuthSignInPage />} />
      <Route path="/lender" element={<Lender />} />
      <Route path="/logout" element= {<Logout />} />
    </Routes>
  </BrowserRouter>
  </ApolloProvider>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
