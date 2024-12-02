import { useEffect, useState } from "react";
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  FlatList,
  TextInput,
  Modal,
  Alert,
  TouchableWithoutFeedback,
} from "react-native";
import { Feather, Ionicons } from "@expo/vector-icons";
import { router, useLocalSearchParams } from "expo-router";
import axios from "axios";
import { useGlobalContext } from "@/context/GlobalProvider";
import { chatServiceHost } from "@/constants/backendUrl";

const ParticipantsPage = () => {
  const [participants, setParticipants] = useState([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [isAddModalVisible, setAddModalVisible] = useState(false);
  const [newParticipant, setNewParticipant] = useState("");
  const [availableUsers, setAvailableUsers] = useState([]); // List of all users
  const [selectedUsers, setSelectedUsers] = useState<any>([]);

  const { chatId } = useLocalSearchParams();
  const context = useGlobalContext();
  if (context == undefined) throw new Error("Context not defined");

  const { token } = context;

  const toggleUserSelection = (userId: any) => {
    setSelectedUsers((prev: any) =>
      prev.includes(userId)
        ? prev.filter((id: any) => id !== userId)
        : [...prev, userId]
    );
  };

  const handleAddSelectedParticipants = async () => {
    try {
      await axios.post(
        `${chatServiceHost}/${chatId}/addParticipants`,
        { userIds: selectedUsers },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      setSelectedUsers([]);
      setAddModalVisible(false);
      fetchParticipants();
    } catch (error) {
      console.error("Failed to add participants", error);
      Alert.alert("Error", "Could not add participants. Please try again.");
    }
  };

  useEffect(() => {
    if (!token) {
      router.replace("/");
      return;
    }
    fetchParticipants();
  }, []);

  const fetchParticipants = async () => {
    try {
      const { data } = await axios.get(`${chatServiceHost}/getChatMembers`, {
        headers: { Authorization: `Bearer ${token}`, chat_id: chatId },
      });
      setParticipants(data);
      console.log(data);
    } catch (error) {
      console.error("Failed to fetch participants", error);
    }
  };

  const addParticipant = async () => {
    if (!newParticipant.trim()) {
      Alert.alert("Error", "Participant ID cannot be empty.");
      return;
    }
    try {
      await axios.post(
        `${chatServiceHost}/chat/${chatId}/addParticipant`,
        { userId: newParticipant },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      setNewParticipant("");
      setAddModalVisible(false);
      //   fetchParticipants();
    } catch (error) {
      console.error("Failed to add participant", error);
      Alert.alert("Error", "Could not add participant. Please try again.");
    }
  };

  const renderParticipant = ({ item }: { item: any }) => (
    <View style={styles.participantItem}>
      <Text style={styles.participantName}>{item.name}</Text>
    </View>
  );

  const filteredUsers = availableUsers.filter((user: any) =>
    user.name.toLowerCase().includes(searchQuery.toLowerCase())
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
        keyExtractor={item => item.id.toString()}
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
              keyExtractor={item => item.id.toString()}
              renderItem={({ item }: { item: any }) => (
                <TouchableOpacity
                  style={styles.userItem}
                  onPress={() => toggleUserSelection(item.id)}
                >
                  <Text style={styles.userName}>{item.name}</Text>
                  <Ionicons
                    name={
                      selectedUsers.includes(item.id)
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
