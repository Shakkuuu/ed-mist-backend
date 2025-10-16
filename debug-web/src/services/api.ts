import axios from 'axios';
import { Organization, User, Room, Device, Subject, Lesson } from '../types';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8081';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// GET methods
export const getOrganizations = async (): Promise<Organization[]> => {
  const response = await api.get('/debug/organizations');
  return response.data as Organization[];
};

export const getUsers = async (): Promise<User[]> => {
  const response = await api.get('/debug/users');
  return response.data as User[];
};

export const getRooms = async (): Promise<Room[]> => {
  const response = await api.get('/debug/rooms');
  return response.data as Room[];
};

export const getDevices = async (): Promise<Device[]> => {
  const response = await api.get('/debug/devices');
  return response.data as Device[];
};

export const getSubjects = async (): Promise<Subject[]> => {
  const response = await api.get('/debug/subjects');
  return response.data as Subject[];
};

export const getLessons = async (): Promise<Lesson[]> => {
  const response = await api.get('/debug/lessons');
  return response.data as Lesson[];
};

// POST methods
export const createOrganization = async (data: any): Promise<Organization> => {
  const response = await api.post('/debug/organizations', data);
  return response.data as Organization;
};

export const createUser = async (data: any): Promise<User> => {
  const response = await api.post('/debug/users', data);
  return response.data as User;
};

export const createRoom = async (data: any): Promise<Room> => {
  const response = await api.post('/debug/rooms', data);
  return response.data as Room;
};

export const createDevice = async (data: any): Promise<Device> => {
  const response = await api.post('/debug/devices', data);
  return response.data as Device;
};

export const createSubject = async (data: any): Promise<Subject> => {
  const response = await api.post('/debug/subjects', data);
  return response.data as Subject;
};

export const createLesson = async (data: any): Promise<Lesson> => {
  const response = await api.post('/debug/lessons', data);
  return response.data as Lesson;
};

// DELETE methods
export const deleteOrganizations = async (): Promise<{ message: string }> => {
  const response = await api.delete('/debug/organizations');
  return response.data as { message: string };
};

export const deleteUsers = async (): Promise<{ message: string }> => {
  const response = await api.delete('/debug/users');
  return response.data as { message: string };
};

export const deleteRooms = async (): Promise<{ message: string }> => {
  const response = await api.delete('/debug/rooms');
  return response.data as { message: string };
};

export const deleteDevices = async (): Promise<{ message: string }> => {
  const response = await api.delete('/debug/devices');
  return response.data as { message: string };
};

export const deleteSubjects = async (): Promise<{ message: string }> => {
  const response = await api.delete('/debug/subjects');
  return response.data as { message: string };
};

export const deleteLessons = async (): Promise<{ message: string }> => {
  const response = await api.delete('/debug/lessons');
  return response.data as { message: string };
};

// Special methods
export const createSeedData = async (): Promise<{ message: string }> => {
  const response = await api.post('/debug/seed');
  return response.data as { message: string };
};

export const resetDatabase = async (): Promise<{ message: string }> => {
  const response = await api.delete('/debug/reset');
  return response.data as { message: string };
};
