import { useState, useEffect } from "react"
import { StockTicker } from "./StockTicker";
import { useTransport } from "./transport";
import './App.css';

function App() {
  
  const bus = useTransport();
  const [connected, setConnected] = useState(false);

  // connect to broker if we're not connected
  useEffect(
    () => {
      if (!connected) {
        bus.fabric.connect(
          () => {
            // handle success
            console.log("application has connected to the broker");
            setConnected(true);
          },
          () => {
            // handle failure
            console.log("application has disconnected from the broker");
            setConnected(false);
          },
          "transport-bus.io",
          443,
          "/ws",
          true,
          "/topic",
          "/queue"
        );
      }
    }
  );

  
  return (
    <div className="App">
        <h1>Stock Lookup</h1>
        <StockTicker />
    </div>
  );
}

export default App;
