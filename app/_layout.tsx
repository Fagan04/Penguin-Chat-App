import { Stack } from "expo-router";
import FlashMessage from "react-native-flash-message";

const Layout = () => {
  return (
    <>
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
      </Stack>
      <FlashMessage position="top" />
    </>
  );
};

export default Layout;
