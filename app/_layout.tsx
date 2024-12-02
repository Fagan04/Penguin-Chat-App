import GlobalProvider from "@/context/GlobalProvider";
import { Stack } from "expo-router";
import FlashMessage from "react-native-flash-message";

const Layout = () => {
  return (
    <GlobalProvider>
      <Stack
        screenOptions={{
          headerStyle: {
            backgroundColor: "#E5F4FD", // Light blue header background
          },
          headerTintColor: "#000", // Header text/icon color
          headerTitleStyle: {
            fontWeight: "bold",
          },
        }}
      >
        <Stack.Screen name="index" options={{ headerShown: false }} />
        <Stack.Screen name="login" options={{ headerShown: false }} />
        <Stack.Screen
          name="chats/participants/[chatId]"
          options={{ headerShown: false }}
        />
        <Stack.Screen name="chats/index" options={{ headerShown: false }} />
        <Stack.Screen name="chats/new" options={{ headerShown: false }} />
        <Stack.Screen name="chats/[chatId]" options={{ headerShown: false }} />
      </Stack>
      <FlashMessage position="top" />
    </GlobalProvider>
  );
};

export default Layout;
