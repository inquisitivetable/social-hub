import useWebSocket from "react-use-websocket";

const useWebSocketConnection = (socketUrl) => {
  const { sendJsonMessage, lastJsonMessage } = useWebSocket(socketUrl, {
    share: true,
  });

  return { sendJsonMessage, lastJsonMessage };
};

export default useWebSocketConnection;
