// Navigation types for React Navigation
export type RootStackParamList = {
    Login: undefined;
    ChatList: undefined;
    Chat: { chatId: string; chatName: string };
    Profile: undefined;
};