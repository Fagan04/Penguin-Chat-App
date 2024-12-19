import {
  createContext,
  ReactNode,
  useContext,
  useEffect,
  useState,
} from "react";

import AsyncStorage from "@react-native-async-storage/async-storage";
import { Chat } from "@/types/Chat";

interface GlobalContextProps {
  token: string;
  setToken: (token: string) => void;
  chats: Chat[];
  setChats: (chats: Chat[]) => void;
  currentChat: Chat | undefined;
  setCurrentChat: (chats: Chat) => void;
  connection: WebSocket | undefined;
  setConnection: (connection: WebSocket) => void;
}

type Message = {
  message: string;
  user_id: number;
};

const GlobalContext = createContext<GlobalContextProps | undefined>(undefined);
export const useGlobalContext = () => useContext(GlobalContext);

const GlobalProvider = ({ children }: { children: ReactNode }) => {
  const [token, setToken] = useState("");
  const [chats, setChats] = useState<Chat[]>([]);
  const [currentChat, setCurrentChat] = useState<Chat>();
  const [connection, setConnection] = useState<WebSocket>();

  useEffect(() => {
    const getData = async () => {
      const data = localStorage.getItem("token");
      if (data) setToken(JSON.parse(data));
    };
    getData();
  }, []);

  return (
    <GlobalContext.Provider
      value={{
        token,
        setToken,
        chats,
        setChats,
        currentChat,
        setCurrentChat,
        connection,
        setConnection,
      }}
    >
      {children}
    </GlobalContext.Provider>
  );
};

export default GlobalProvider;
