import { createContext, ReactNode, useContext, useState } from "react";

interface GlobalContextProps {
  user: User;
  isLoggedIn: boolean;
  setUser: (user: User) => void;
  setIsLoggedIn: (status: boolean) => void;
  loading: boolean;
}

type User = {
  id: number | null;
  name: string | null;
  email: string | null;
};

const GlobalContext = createContext<GlobalContextProps | undefined>(undefined);
export const useGlobalContext = () => useContext(GlobalContext);

const GlobalProvider = ({ children }: { children: ReactNode }) => {
  const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false);
  const [user, setUser] = useState<User>({ id: null, name: null, email: null });
  const [loading, setLoading] = useState(true);

  return (
    <GlobalContext.Provider
      value={{
        user,
        setIsLoggedIn,
        isLoggedIn,
        setUser,
        loading,
      }}
    >
      {children}
    </GlobalContext.Provider>
  );
};

export default GlobalProvider;
