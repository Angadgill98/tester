import { create } from "zustand";
import { Action_f, Application_f, ChildApplicationInfo, Environment_f, Sequence_f, Tab } from "./models/application";

interface GlobalStore {
    Applications: Map<number, Application_f>;
    Actions: Map<number, Action_f>;
    Sequences: Map<number, Sequence_f>;
    Environments: Map<number, Environment_f>;

    Tabs: Tab[];
    CurrentTab: Tab | null;

    CurrentApplicationId: number | null;

    setApplications: (id: number, app: Application_f) => void;
    setActions: (id: number, action: Action_f) => void;
    setSequences: (id: number, sequence: Sequence_f) => void;
    setEnvironments: (id: number, environment: Environment_f) => void;

    setTabs: (tabs: Tab[]) => void;
    setCurrentTab: (tab: Tab | null) => void;

    setCurrentApplicationId: (id: number | null) => void;
    getCurrentApplicationId: () => number | null;
    getCurrentApplication: () => Application_f | undefined;

    addToTabs: (tab: Tab) => void;
    removeFromTabs: (tab: Tab) => void;
}

export const useGlobalStore = create<GlobalStore>((set, get) => ({
    Applications: new Map<number, Application_f>(),
    Actions: new Map<number, Action_f>(),
    Sequences: new Map<number, Sequence_f>(),
    Environments: new Map<number, Environment_f>(),
    ChildApplications: new Map<number, ChildApplicationInfo>(),

    Tabs: [],
    CurrentTab: null,

    CurrentApplicationId: null,

    setApplications: (id, app) => set((state) => {
        const applications = new Map(state.Applications);
        applications.set(id, app);
        return { Applications: applications };
    }),

    setActions: (id, action) => set((state) => {
        const actions = new Map(state.Actions);
        actions.set(id, action);
        return { Actions: actions };
    }),

    setSequences: (id, sequence) => set((state) => {
        const sequences = new Map(state.Sequences);
        sequences.set(id, sequence);
        return { Sequences: sequences };
    }),

    setEnvironments: (id, environment) => set((state) => {
        const environments = new Map(state.Environments);
        environments.set(id, environment);
        return { Environments: environments };
    }),


    setTabs: (tabs) => set({ Tabs: tabs }),

    setCurrentTab: (tab) => set({ CurrentTab: tab }),


    setCurrentApplicationId: (id) => set({ CurrentApplicationId: id }),

    getCurrentApplicationId: () => get().CurrentApplicationId,

    getCurrentApplication: () => {
        const state = get();

        if (state.CurrentApplicationId === null) return undefined;

        return state.Applications.get(state.CurrentApplicationId);
    },



    addToTabs: (tab) => set((state) => {
        const existingTab = state.Tabs.find(
            t => t.TypeId === tab.TypeId &&
                t.AppId === tab.AppId &&
                t.TabType === tab.TabType
        );

        const tabs = existingTab ? state.Tabs : [...state.Tabs, tab];

        return {
            Tabs: tabs,
            CurrentTab: existingTab ?? tab
        };
    }),


    removeFromTabs: (tab) => set((state) => {
        const tabs = state.Tabs.filter(
            t => !(
                t.TypeId === tab.TypeId &&
                t.AppId === tab.AppId &&
                t.TabType === tab.TabType
            )
        );

        return {
            Tabs: tabs,
            CurrentTab:
                state.CurrentTab &&
                state.CurrentTab.TypeId === tab.TypeId &&
                state.CurrentTab.AppId === tab.AppId &&
                state.CurrentTab.TabType === tab.TabType
                    ? null
                    : state.CurrentTab
        };
    }),



}));