import { createContext } from 'react';
import { DomainUser, DomainAdminUser } from '@/api/types';

export const AuthContext = createContext<
  [
    DomainUser | DomainAdminUser | null,
    {
      loading: boolean;
      setUser: (user: DomainUser | DomainAdminUser) => void;
      refreshUser: () => void;
    }
  ]
>([
  null,
  {
    setUser: () => {},
    loading: true,
    refreshUser: () => {},
  },
]);

export const CommonContext = createContext<{
  contactModalOpen: boolean;
  setContactModalOpen: React.Dispatch<React.SetStateAction<boolean>>;
}>({ contactModalOpen: false, setContactModalOpen: () => {} });
