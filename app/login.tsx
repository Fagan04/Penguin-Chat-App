import { userServiceHost } from "@/constants/backendUrl";
import { useGlobalContext } from "@/context/GlobalProvider";
import showSuccessMessage from "@/utils/showSuccessMessage";
import axios from "axios";
import { Link, useRouter } from "expo-router";
import { useEffect, useState } from "react";
import {
  View,
  Text,
  TextInput,
  StyleSheet,
  ImageBackground,
  Pressable,
} from "react-native";

const image = {
  uri: "https://s3-alpha-sig.figma.com/img/745d/f541/dcc60ddf9129d621d45f365b41b0dfb8?Expires=1732492800&Key-Pair-Id=APKAQ4GOSFWCVNEHN3O4&Signature=jM~PSzgNAyAYmhqQSeg9llZnFtdTH4etWpTETGbvnH6QsRXrNc9DRZ3aAxmlff4iLFDwYUC16P57NKqXs1kPuxLqISZw6SESqVi6HtXWt6esOGPDRWpWanNRNnLQo29b5oirGHFJzWHbrNKtGGoRmNLui0reXesCsuafCwTnbxv4yxKHxf3nEJNfgDagW7sRMaJptcs9rMLrr-JedcJDhcCLagcLgj4-0MTmTsONzbHnbqsIK78qfDTDBxS33xRaOaFNgePdYFzouV697B0UGIcgpLamgofZb3q0hmZlZwMwa-IlU-OMXrXWdE3hacJ3u0V55O8qyLazdLRYFkI6ow__",
};

const Login = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const context = useGlobalContext();
  const router = useRouter();

  if (context == undefined) throw new Error("Context not defined");

  const { setToken, token } = context;

  useEffect(() => {
    if (token) router.replace("/chats");
  }, [token]);

  const handleSubmit = async () => {
    try {
      const { data }: { data: { message: string; token: string } } =
        await axios.post(`${userServiceHost}/login`, {
          username,
          password,
        });
      showSuccessMessage(data.message);
      setToken(data.token);
      router.replace("/chats");
    } catch (error) {
      if (axios.isAxiosError(error)) {
        showSuccessMessage(error.response?.data.trim(), false);
      } else {
        console.log("Unexpected error:", error);
      }
    }
  };

  return (
    <ImageBackground
      source={image}
      style={styles.background}
      resizeMode="cover"
    >
      <View style={styles.container}>
        <View style={styles.topContainer}>
          <Text style={styles.title}>Log In</Text>
          <TextInput
            placeholderTextColor={"#e5e5e5"}
            placeholder="Username"
            style={styles.input}
            value={username}
            onChangeText={e => setUsername(e)}
          />
          <TextInput
            secureTextEntry
            placeholderTextColor={"#e5e5e5"}
            placeholder="Password"
            style={styles.input}
            value={password}
            onChangeText={e => setPassword(e)}
          />
          <Pressable
            onPress={handleSubmit}
            style={({ pressed }) => [
              {
                backgroundColor: pressed ? "#007aff50" : "#007aff",
              },
              styles.submit,
            ]}
          >
            <Text style={styles.submitText}>Log In</Text>
          </Pressable>
          <Text style={styles.info}>
            New to LosPenguinos?{" "}
            <Link replace href="/" style={styles.infoLink}>
              Sign up here
            </Link>
          </Text>
        </View>
        <View style={styles.bottomContainer}>
          <Text style={{ fontSize: 32, fontWeight: 700 }}>PenguinChat</Text>
        </View>
      </View>
    </ImageBackground>
  );
};

const styles = StyleSheet.create({
  background: {
    flex: 1,
    justifyContent: "center",
    width: "100%",
    height: "100%",
  },
  container: {
    width: "100%",
    flexDirection: "column",
    height: "80%",
    alignItems: "center",
  },
  title: {
    fontSize: 32,
    fontWeight: 700,
    color: "#007AFF",
    marginBottom: 30,
    textAlign: "center",
  },
  topContainer: {
    flexDirection: "column",
    width: "100%",
    justifyContent: "flex-start",
    alignItems: "center",
    gap: 10,
  },
  bottomContainer: {
    height: "90%",
    justifyContent: "center",
  },
  input: {
    marginBottom: 15,
    backgroundColor: "white",
    paddingInline: 35,
    paddingBlock: 15,
    borderRadius: 5,
    width: "65%",
    position: "relative",
  },
  submit: {
    width: "65%",
    marginTop: 10,
    paddingInline: 35,
    paddingBlock: 15,
    borderRadius: 5,
  },
  submitText: {
    textAlign: "center",
    color: "white",
    fontSize: 18,
  },
  info: {
    marginTop: 5,
  },
  infoLink: {
    color: "#007aff",
  },
});

export default Login;
