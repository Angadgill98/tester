import { createContext, useState, type ReactNode } from "react";
import { Application_f } from "./models/application";

interface GlobalContextstruct {
    GlobalApplication: Application_f | null;
    setGlobalApplication: React.Dispatch<React.SetStateAction<Application_f | null>>;

    CurrentApplication: Application_f | null;
    setCurrentApplication: React.Dispatch<React.SetStateAction<Application_f | null>>;
}

export const GlobalContext = createContext<GlobalContextstruct | null>(null);

interface GlobalProviderProps {
    children: ReactNode;
}

export function GlobalProvider({ children }: GlobalProviderProps) {
    const [GlobalApplication, setGlobalApplication] = useState<Application_f | null>(null);
    const [CurrentApplication, setCurrentApplication] = useState<Application_f | null>(null);

    return (
        <GlobalContext.Provider value={{
            GlobalApplication,
            setGlobalApplication,
            CurrentApplication,
            setCurrentApplication
        }}>
            {children}
        </GlobalContext.Provider>
    );
}