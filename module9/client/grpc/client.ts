import { createClient, type Client } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { type DescService } from "@bufbuild/protobuf";
import { useMemo } from "react";

const transport = createConnectTransport({
    baseUrl: "http://localhost:50052",
    useBinaryFormat: false, // optional, JSON easier for debugging
    // Not needed. Web browsers use HTTP/2 automatically.
    // httpVersion: "1.1"
});

export function useClient<T extends DescService>(service: T): Client<T> {
    // We memoize the client, so that we only create one instance per service.
    return useMemo(() => createClient(service, transport), [service]);
}