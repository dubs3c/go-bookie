
import { writable } from 'svelte/store';


interface Settings {
  currentPage: number
  totalPages: number
  pageSize: number
  archiveChecked: boolean
  deletedChecked: boolean
  activeTags: Record<string, string>
};

const initialState: Settings = {
  currentPage: 1,
  totalPages: 1,
  pageSize: 10,
  archiveChecked: false,
  deletedChecked: false,
  activeTags: {}
}

const AuthLocalStore = () => {
  return {
    addToken: (token: string) =>
      localStorage.setItem("accesstoken", token),
    getToken: () =>
      localStorage.getItem("accesstoken"),
    removeToken: () =>
      localStorage.removeItem("accesstoken")
  }
}

const PageSettings = () => {
  // creates a new writable store populated with some initial data
  const { subscribe, update, set } = writable(initialState);

  return {
    subscribe,
    update,
    set,
    setCurrentPage: pageNumber => 
      update(state => (state = {...state, currentPage: pageNumber})),
    setTotalPages: amountPages =>
      update(state => (state = {...state, totalPages: amountPages})),
    setPageSize: pageSize =>
      update(state => (state = {...state, pageSize: pageSize})),
    setArchived: checked =>
      update(state => (state = {...state, archiveChecked: checked})),
    setDeleted: checked =>
      update(state => (state = {...state, deletedChecked: checked})),
    setActiveTags: tags =>
      update(state => (state = {...state, activeTags: tags})),
    setAll: all =>
      update(state => (state = {...all})),     
  };
};

export const authStore = AuthLocalStore()
export const settingsStore = PageSettings()
export const bookmarkStore = writable([])
export const tagStore = writable([])

// store to handle the app state
const appState = () => {
  const { subscribe, update } = writable(false);
  return {
    subscribe,
    error: () => update(error => !error),
  };
};

export const AppStore = appState();