import React, { useState } from "react";
import { render, fireEvent } from "@testing-library/react-native";
import { Text, TextInput, Pressable, View } from "react-native";

const Login = () => {
  const [successMessage, setSuccessMessage] = useState("");

  const handleSubmit = () => {
    setSuccessMessage("Login successful");
  };

  return (
    <View>
      <TextInput placeholder="Username" testID="username" />
      <TextInput placeholder="Password" testID="password" secureTextEntry />
      <Pressable onPress={handleSubmit} testID="submit">
        <Text>Log In</Text>
      </Pressable>
      {successMessage ? <Text testID="successMessage">{successMessage}</Text> : null}
    </View>
  );
};

describe("Login Component", () => {
  it("should display success message on login", () => {
    const { getByTestId, queryByTestId } = render(<Login />);

    // Ensure the success message is not shown initially
    expect(queryByTestId("successMessage")).toBeNull();

    // Simulate login button press
    const submitButton = getByTestId("submit");
    fireEvent.press(submitButton);

    // Check if the success message is displayed
    expect(getByTestId("successMessage").props.children).toBe("Login successful");
  });
});

