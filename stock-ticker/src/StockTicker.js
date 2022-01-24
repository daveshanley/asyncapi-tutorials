import { useState, useEffect } from "react";
import { useTransport } from "./transport";
import { GeneralUtil } from "@vmw/transport/util/util";

// stock ticker component
export function StockTicker() {

    // define hooks
    const [stockLookup, setStockLookup] = useState("GOOG");
    const [closePrice, setClosePrice] = useState(null);
    const [lastRefreshed, setLastRefreshed] = useState(null);
    const [symbol, setSymbol] = useState(null);
    const [stockError, setStockError] = useState(null);
    const [subscribed, setSubscribed] = useState(false);

    // grab a reference to our event bus
    const bus = useTransport();
    
    // define the channel our stock service is operating on.

    const stockChannel = "stock-ticker-service";

    // create a connection string to identify our broker
    const connectionString = GeneralUtil.getFabricConnectionString(
        "transport-bus.io",
        443,
        "/ws"
    );

    // if we are subscribed, then we do nothing, if not... we subscribe!
    useEffect(
        () => {
            if (!subscribed) {
                // subscribe to service by marking the channel as galactic.
                bus.markChannelAsGalactic(stockChannel, connectionString, true);
                setSubscribed(true);
            }
            return () => {
                // all done, unsubscribe
                bus.markChannelAsLocal(stockChannel);
            };
        }, [subscribed, bus, connectionString]
    );

    function handleSymbolChange(e) {
        setStockLookup(e.target.value);
    }

    function requestStock() {
        let request = bus.fabric.generateFabricRequest(
            "ticker_price_update_stream",
            { symbol: stockLookup }
        );

        // make the request over the bus and over to our broker at transport-bus.io
        bus.requestOnce(stockChannel, request).handle(
            (response) => {
                // success
                if (response.payload) {
                    setClosePrice(response.payload.closePrice);
                    setLastRefreshed(response.payload.lastRefreshed);
                    setSymbol(response.payload.symbol);
                } else {
                    setStockError("nothing returned by service");
                }
            },
            (error) => {
                // error
                setStockError(error.errorMessage);
            }
        );
    }

    // return our component JSX
    return (
        <div>
            <label>Ticker Symbol: </label>
            <input type="text" value={stockLookup} onChange={handleSymbolChange} />
            <button onClick={requestStock}>Get Price!</button>

            <StockError errorMessage={stockError} />
            <StockResult
                closePrice={closePrice}
                symbol={symbol}
                lastRefreshed={lastRefreshed}
            />
        </div>
    );
}

function StockError(props) {
    if (props.errorMessage) {
        return (
            <div>
                <hr />
                <h3>Sorry! the service issued an error</h3>
                <p> {props.errorMessage} </p>
            </div>
        );
    }
    return null;
}


function StockResult(props) {
    let price = props.closePrice?.toFixed(2);
    let lastRefreshed = props.lastRefreshed;
    let symbol = props.symbol;

    if (price && lastRefreshed && symbol) {
        return (
            <div>
                <hr/>
                Symbol: <strong>{symbol}</strong>
                <br/>
                Price: <strong>{price}</strong>
                <br/>
                Last Refreshed: <strong>{lastRefreshed}</strong>
            </div>
        );
    }
    return null;
}