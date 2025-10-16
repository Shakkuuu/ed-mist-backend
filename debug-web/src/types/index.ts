export interface Organization {
  id: string;
  mail: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface User {
  id: string;
  org_id: string;
  mail: string;
  created_at: string;
  updated_at: string;
}

export interface Room {
  id: string;
  org_id: string;
  org_room_id: string;
  name: string;
  caption: string;
  mist_zone_id: string;
  created_at: string;
  updated_at: string;
}

export interface Device {
  id: string;
  user_id: string;
  device_id: string;
  is_active: boolean;
  last_authenticated: string;
  created_at: string;
  updated_at: string;
}

export interface Subject {
  id: string;
  org_id: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface Lesson {
  id: string;
  subject_id: string;
  room_id: string;
  org_id: string;
  day_of_week: number;
  start_time: string;
  end_time: string;
  period: number;
  created_at: string;
  updated_at: string;
}

export type ResourceType = 'organizations' | 'users' | 'rooms' | 'devices' | 'subjects' | 'lessons';
