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
}

type User = {
  id: number | null;
  name: string | null;
  email: string | null;
};

const GlobalContext = createContext<GlobalContextProps | undefined>(undefined);
export const useGlobalContext = () => useContext(GlobalContext);

const GlobalProvider = ({ children }: { children: ReactNode }) => {
  const [token, setToken] = useState("");
  const [chats, setChats] = useState<Chat[]>([]);
  const [currentChat, setCurrentChat] = useState<Chat>();

  useEffect(() => {
    const getData = async () => {
      const data = await AsyncStorage.getItem("token");
      if (data) setToken(data);
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
      }}
    >
      {children}
    </GlobalContext.Provider>
  );
};

export default GlobalProvider;
