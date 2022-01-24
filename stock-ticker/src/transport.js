import { BusUtil } from "@vmw/transport/util/bus.util"

export function useTransport() {
    return BusUtil.getBusInstance();
}