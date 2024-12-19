import { useEffect, useState } from "react";
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  FlatList,
  TextInput,
  Modal,
} from "react-native";
import { Feather, Ionicons } from "@expo/vector-icons";
import { router, useLocalSearchParams } from "expo-router";
import axios from "axios";
import { useGlobalContext } from "@/context/GlobalProvider";
import { chatServiceHost } from "@/constants/backendUrl";
import { Participant } from "@/types/Participant";
import showSuccessMessage from "@/utils/showSuccessMessage";
import { Owner } from "@/types/Owner";
import { jwtDecode } from "jwt-decode";

const ParticipantsPage = () => {
  const [participants, setParticipants] = useState<Participant[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [isAddModalVisible, setAddModalVisible] = useState(false);
  const [availableUsers, setAvailableUsers] = useState<string[]>([]); // List of all users
  const [selectedUsers, setSelectedUsers] = useState<string[]>([]);
  const [owner, setOwner] = useState<Owner>();

  const { chatId } = useLocalSearchParams();

  const context = useGlobalContext();
  if (context == undefined) throw new Error("Context not defined");

  const { token } = context;
  const { username: currentUsername }: { username: string } = jwtDecode(token);

  const toggleUserSelection = (userId: string) => {
    setSelectedUsers((prev: string[]) =>
      prev.includes(userId)
        ? prev.filter((id: string) => id !== userId)
        : [...prev, userId]
    );
  };

  const handleAddSelectedParticipants = async () => {
    try {
      for (const username of selectedUsers) {
        await axios.post(
          `${chatServiceHost}/addUserToChat`,
          { username: username, chat_id: Number(chatId) },
          { headers: { Authorization: `Bearer ${token}` } }
        );
      }
      setSelectedUsers([]);
      setAddModalVisible(false);
      fetchParticipants();
    } catch (error) {
      console.error("Failed to add participants", error);
      showSuccessMessage("Failed to add participant(s).", false);
    }
  };

  useEffect(() => {
    if (!token) {
      router.replace("/");
      return;
    }
    fetchParticipants();
    fetchAllUsers();
  }, []);

  const fetchAllUsers = async () => {
    let { data }: { data: string[] } = await axios.get(
      `${chatServiceHost}/getAllUsers`,
      {
        headers: { Authorization: `Bearer ${token}` },
      }
    );
    data = data.filter(username => username !== currentUsername);
    setAvailableUsers(data);
  };

  const fetchParticipants = async () => {
    try {
      const {
        data,
      }: {
        data: {
          owner_id: number;
          owner_username: string;
          participants: Participant[];
        };
      } = await axios.get(`${chatServiceHost}/getChatMembers`, {
        headers: { Authorization: `Bearer ${token}`, chat_id: chatId },
      });
      setOwner({
        owner_id: data.owner_id,
        owner_username: data.owner_username,
      });
      setParticipants(data.participants);
    } catch (error) {
      console.error("Failed to fetch participants", error);
    }
  };

  const renderParticipant = ({ item }: { item: Participant }) => (
    <View style={styles.participantItem}>
      <Text style={styles.participantName}>{item.Username}</Text>
      {item.UserID === owner?.owner_id ? (
        <Text style={styles.participantName}>Owner</Text>
      ) : (
        ""
      )}
    </View>
  );

  const filteredUsers = availableUsers.filter((user: string) =>
    user.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <TouchableOpacity onPress={router.back}>
          <Feather name="arrow-left" size={24} color="black" />
        </TouchableOpacity>
        <Text style={styles.headerTitle}>Participants</Text>
        <TouchableOpacity onPress={() => setAddModalVisible(true)}>
          <Ionicons name="add" size={24} color="black" />
        </TouchableOpacity>
      </View>

      <FlatList
        data={participants}
        keyExtractor={item => item.UserID.toString()}
        renderItem={renderParticipant}
        contentContainerStyle={styles.participantsList}
      />

      <Modal
        animationType="slide"
        transparent={true}
        visible={isAddModalVisible}
        onRequestClose={() => setAddModalVisible(false)}
      >
        <View style={styles.modalOverlay}>
          <View style={styles.modalContent}>
            <Text style={styles.modalTitle}>Add Participants</Text>
            <TextInput
              style={styles.input}
              placeholder="Search users"
              value={searchQuery}
              onChangeText={setSearchQuery}
            />
            <FlatList
              data={filteredUsers} // Assuming `availableUsers` is a list of users fetched from the server
              keyExtractor={item => item}
              renderItem={({ item }: { item: string }) => (
                <TouchableOpacity
                  style={styles.userItem}
                  onPress={() => toggleUserSelection(item)}
                >
                  <Text style={styles.userName}>{item}</Text>
                  <Ionicons
                    name={
                      selectedUsers.includes(item)
                        ? "checkbox"
                        : "square-outline"
                    }
                    size={24}
                    color="#2b59c3"
                  />
                </TouchableOpacity>
              )}
              contentContainerStyle={styles.userList}
            />
            <TouchableOpacity
              style={styles.addButton}
              onPress={handleAddSelectedParticipants}
            >
              <Text style={styles.addButtonText}>Done</Text>
            </TouchableOpacity>
            <TouchableOpacity
              style={styles.closeButton}
              onPress={() => setAddModalVisible(false)}
            >
              <Text style={styles.closeButtonText}>Cancel</Text>
            </TouchableOpacity>
          </View>
        </View>
      </Modal>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#f9f9f9",
  },
  header: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
    padding: 16,
    paddingTop: 40,
    backgroundColor: "#ffffff",
    borderBottomWidth: 1,
    borderBottomColor: "#ddd",
  },
  headerTitle: {
    fontSize: 18,
    fontWeight: "bold",
  },
  participantsList: {
    padding: 16,
  },
  participantItem: {
    display: "flex",
    flexDirection: "row",
    justifyContent: "space-between",
    padding: 12,
    borderBottomWidth: 1,
    borderBottomColor: "#ddd",
  },
  participantName: {
    fontSize: 16,
  },
  modalOverlay: {
    flex: 1,
    backgroundColor: "rgba(0, 0, 0, 0.5)",
    justifyContent: "center",
    alignItems: "center",
  },
  modalContent: {
    width: "90%",
    backgroundColor: "white",
    borderRadius: 10,
    padding: 20,
    alignItems: "center",
  },
  modalTitle: {
    fontSize: 20,
    fontWeight: "bold",
    marginBottom: 10,
  },
  input: {
    width: "100%",
    height: 40,
    borderWidth: 1,
    borderColor: "#ddd",
    borderRadius: 5,
    paddingHorizontal: 10,
    marginBottom: 20,
  },
  addButton: {
    backgroundColor: "#2b59c3",
    paddingVertical: 10,
    paddingHorizontal: 20,
    borderRadius: 5,
    marginBottom: 10,
  },
  addButtonText: {
    color: "white",
    fontSize: 16,
    fontWeight: "bold",
  },
  closeButton: {
    backgroundColor: "#f44336",
    paddingVertical: 10,
    paddingHorizontal: 20,
    borderRadius: 5,
  },
  closeButtonText: {
    color: "white",
    fontSize: 16,
    fontWeight: "bold",
  },
  userList: {
    width: "100%",
    maxHeight: 300,
    marginBottom: 20,
  },
  userItem: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    paddingVertical: 10,
    paddingHorizontal: 15,
    borderBottomWidth: 1,
    borderBottomColor: "#ddd",
  },
  userName: {
    fontSize: 16,
    color: "#333",
  },
});

export default ParticipantsPage;
