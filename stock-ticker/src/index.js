import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';
import { BusUtil } from "@vmw/transport/util/bus.util"
import { LogLevel } from "@vmw/transport/log/logger.model"

ReactDOM.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
  document.getElementById('root')
);

// boot bus and set logging levels.
const bus = BusUtil.bootBus();
bus.api.setLogLevel(LogLevel.Verbose);
bus.log.useDarkTheme(true); // disable if you're not using dark mode in your browser.

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
