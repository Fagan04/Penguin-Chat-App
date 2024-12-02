import { showMessage } from "react-native-flash-message";

export default (
  message: string,
  isSuccess: boolean = true,
  header: string = ""
) => {
  return showMessage({
    message: header ? header : isSuccess ? "Success" : "Error",
    type: isSuccess ? "success" : "danger",
    color: "#fff",
    description: message,
    backgroundColor: "#007AFF",
  });
};
