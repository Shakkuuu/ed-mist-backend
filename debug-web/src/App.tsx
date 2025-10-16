import React, { useState, useEffect, useCallback } from 'react';
import {
  getOrganizations, getUsers, getRooms, getDevices, getSubjects, getLessons,
  createOrganization, createUser, createRoom, createDevice, createSubject, createLesson,
  deleteOrganizations, deleteUsers, deleteRooms, deleteDevices, deleteSubjects, deleteLessons,
  createSeedData, resetDatabase
} from './services/api';
import { Organization, User, Room, Device, Subject, Lesson, ResourceType } from './types';

interface NotificationState {
  show: boolean;
  message: string;
  type: 'success' | 'error' | 'info';
}

const App: React.FC = () => {
  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [rooms, setRooms] = useState<Room[]>([]);
  const [devices, setDevices] = useState<Device[]>([]);
  const [subjects, setSubjects] = useState<Subject[]>([]);
  const [lessons, setLessons] = useState<Lesson[]>([]);

  const [loading, setLoading] = useState<Record<ResourceType, boolean>>({
    organizations: false,
    users: false,
    rooms: false,
    devices: false,
    subjects: false,
    lessons: false,
  });

  const [notification, setNotification] = useState<NotificationState>({
    show: false,
    message: '',
    type: 'info'
  });

  const [formData, setFormData] = useState<Record<ResourceType, any>>({
    organizations: {},
    users: {},
    rooms: {},
    devices: {},
    subjects: {},
    lessons: {},
  });

  const showNotification = useCallback((message: string, type: 'success' | 'error' | 'info') => {
    setNotification({ show: true, message, type });
    setTimeout(() => {
      setNotification(prev => ({ ...prev, show: false }));
    }, 3000);
  }, []);

  const loadData = useCallback(async (resourceType: ResourceType) => {
    setLoading(prev => ({ ...prev, [resourceType]: true }));
    try {
      switch (resourceType) {
        case 'organizations':
          const orgs = await getOrganizations();
          setOrganizations(orgs);
          break;
        case 'users':
          const users = await getUsers();
          setUsers(users);
          break;
        case 'rooms':
          const rooms = await getRooms();
          setRooms(rooms);
          break;
        case 'devices':
          const devices = await getDevices();
          setDevices(devices);
          break;
        case 'subjects':
          const subjects = await getSubjects();
          setSubjects(subjects);
          break;
        case 'lessons':
          const lessons = await getLessons();
          setLessons(lessons);
          break;
      }
    } catch (error: any) {
      showNotification(`データの読み込みに失敗しました: ${error.message}`, 'error');
    } finally {
      setLoading(prev => ({ ...prev, [resourceType]: false }));
    }
  }, [showNotification]);

  const clearForm = (resourceType: ResourceType) => {
    setFormData(prev => ({
      ...prev,
      [resourceType]: {}
    }));
  };

  const createData = async (resourceType: ResourceType, data: any) => {
    setLoading(prev => ({ ...prev, [resourceType]: true }));
    try {
      switch (resourceType) {
        case 'organizations':
          const org = await createOrganization(data);
          setOrganizations(prev => [...prev, org]);
          break;
        case 'users':
          // フィールド名を変換
          const userData = {
            ...data,
            organization_id: data.organization_id,
          };
          const user = await createUser(userData);
          setUsers(prev => [...prev, user]);
          break;
        case 'rooms':
          // フィールド名を変換
          const roomData = {
            ...data,
            organization_id: data.organization_id,
          };
          const room = await createRoom(roomData);
          setRooms(prev => [...prev, room]);
          break;
        case 'devices':
          const device = await createDevice(data);
          setDevices(prev => [...prev, device]);
          break;
        case 'subjects':
          // フィールド名を変換
          const subjectData = {
            ...data,
            organization_id: data.organization_id,
          };
          const subject = await createSubject(subjectData);
          setSubjects(prev => [...prev, subject]);
          break;
        case 'lessons':
          // フィールド名を変換
          const lessonData = {
            ...data,
            org_id: data.organization_id,
          };
          delete lessonData.organization_id;
          const lesson = await createLesson(lessonData);
          setLessons(lessons => [...lessons, lesson]);
          break;
      }
      showNotification(`${getResourceName(resourceType)}を作成しました`, 'success');
      clearForm(resourceType);
    } catch (error: any) {
      showNotification(`作成に失敗しました: ${error.response?.data?.error || error.message}`, 'error');
    } finally {
      setLoading(prev => ({ ...prev, [resourceType]: false }));
    }
  };

  const deleteData = async (resourceType: ResourceType) => {
    if (!window.confirm(`${getResourceName(resourceType)}を全て削除しますか？`)) {
      return;
    }

    setLoading(prev => ({ ...prev, [resourceType]: true }));
    try {
      switch (resourceType) {
        case 'organizations':
          await deleteOrganizations();
          setOrganizations([]);
          break;
        case 'users':
          await deleteUsers();
          setUsers([]);
          break;
        case 'rooms':
          await deleteRooms();
          setRooms([]);
          break;
        case 'devices':
          await deleteDevices();
          setDevices([]);
          break;
        case 'subjects':
          await deleteSubjects();
          setSubjects([]);
          break;
        case 'lessons':
          await deleteLessons();
          setLessons([]);
          break;
      }
      showNotification(`${getResourceName(resourceType)}を削除しました`, 'success');
    } catch (error: any) {
      showNotification(`削除に失敗しました: ${error.response?.data?.error || error.message}`, 'error');
    } finally {
      setLoading(prev => ({ ...prev, [resourceType]: false }));
    }
  };

  const handleSeedData = async () => {
    if (!window.confirm('シードデータを作成しますか？')) {
      return;
    }

    try {
      await createSeedData();
      showNotification('シードデータを作成しました', 'success');
      // 全データを再読み込み
      await Promise.all([
        loadData('organizations'),
        loadData('users'),
        loadData('rooms'),
        loadData('devices'),
        loadData('subjects'),
        loadData('lessons')
      ]);
    } catch (error: any) {
      showNotification(`シードデータの作成に失敗しました: ${error.response?.data?.error || error.message}`, 'error');
    }
  };

  const handleResetDatabase = async () => {
    if (!window.confirm('データベースをリセットしますか？全てのデータが削除されます。')) {
      return;
    }

    try {
      await resetDatabase();
      showNotification('データベースをリセットしました', 'success');
      // 全データを再読み込み
      await Promise.all([
        loadData('organizations'),
        loadData('users'),
        loadData('rooms'),
        loadData('devices'),
        loadData('subjects'),
        loadData('lessons')
      ]);
    } catch (error: any) {
      showNotification(`データベースのリセットに失敗しました: ${error.response?.data?.error || error.message}`, 'error');
    }
  };

  const getResourceName = (resourceType: ResourceType): string => {
    const names: Record<ResourceType, string> = {
      organizations: '組織',
      users: 'ユーザー',
      rooms: '部屋',
      devices: 'デバイス',
      subjects: '科目',
      lessons: '授業'
    };
    return names[resourceType];
  };

  const getFormFields = (resourceType: ResourceType) => {
    switch (resourceType) {
      case 'organizations':
        return [
          { name: 'name', label: '組織名', type: 'text', required: true },
          { name: 'mail', label: 'メールアドレス', type: 'email', required: true },
        ];
      case 'users':
        return [
          { name: 'email', label: 'メールアドレス', type: 'email', required: true },
          { name: 'organization_id', label: '組織ID', type: 'text', required: true },
        ];
      case 'rooms':
        return [
          { name: 'name', label: '部屋名', type: 'text', required: true },
          { name: 'org_room_id', label: '組織部屋ID', type: 'text', required: true },
          { name: 'caption', label: '説明', type: 'textarea', required: false },
          { name: 'mist_zone_id', label: 'MistゾーンID', type: 'text', required: false },
          { name: 'organization_id', label: '組織ID', type: 'text', required: true },
        ];
      case 'devices':
        return [
          { name: 'user_id', label: 'ユーザーID', type: 'text', required: true },
          { name: 'device_id', label: 'デバイスID', type: 'text', required: true },
        ];
      case 'subjects':
        return [
          { name: 'name', label: '科目名', type: 'text', required: true },
          { name: 'year', label: '年代', type: 'number', required: true },
          { name: 'organization_id', label: '組織ID', type: 'text', required: true },
        ];
      case 'lessons':
        return [
          { name: 'subject_id', label: '科目ID', type: 'text', required: true },
          { name: 'room_id', label: '部屋ID', type: 'text', required: true },
          { name: 'organization_id', label: '組織ID', type: 'text', required: true },
          { name: 'start_time', label: '開始時間', type: 'datetime-local', required: true },
          { name: 'end_time', label: '終了時間', type: 'datetime-local', required: true },
        ];
      default:
        return [];
    }
  };

  const handleInputChange = (resourceType: ResourceType, name: string, value: string | number) => {
    setFormData(prev => ({
      ...prev,
      [resourceType]: {
        ...prev[resourceType],
        [name]: value
      }
    }));
  };

  const handleFormSubmit = (e: React.FormEvent, resourceType: ResourceType) => {
    e.preventDefault();
    createData(resourceType, formData[resourceType]);
  };

  useEffect(() => {
    // 初期データ読み込み
    loadData('organizations');
    loadData('users');
    loadData('rooms');
    loadData('devices');
    loadData('subjects');
    loadData('lessons');
  }, [loadData]);

  const renderDataTable = (title: string, resourceType: ResourceType, data: any[]) => {
    return (
      <div className="card">
        <div className="card-header">
          <h3 className="card-title">{title}</h3>
          <span className="card-count">{data.length}件</span>
        </div>
        <div className="card-body">
          <div className="d-flex gap-2 mb-3">
            <button
              className="btn btn-danger"
              onClick={() => deleteData(resourceType)}
              disabled={loading[resourceType]}
            >
              全削除
            </button>
            <button
              className="btn btn-secondary"
              onClick={() => loadData(resourceType)}
              disabled={loading[resourceType]}
            >
              更新
            </button>
          </div>

          {/* 作成フォーム */}
          <div className="form-section mb-4">
            <h4 className="form-section-title">新規作成</h4>
            <form onSubmit={(e) => handleFormSubmit(e, resourceType)}>
              <div className="form-grid">
                {getFormFields(resourceType).map((field) => (
                  <div key={field.name} className="form-group">
                    <label className="form-label">
                      {field.label}
                      {field.required && <span style={{color: 'red'}}> *</span>}
                    </label>

                    {field.type === 'textarea' ? (
                      <textarea
                        className="form-control"
                        value={formData[resourceType][field.name] || ''}
                        onChange={(e) => handleInputChange(resourceType, field.name, e.target.value)}
                        required={field.required}
                      />
                    ) : (
                      <input
                        type={field.type}
                        className="form-control"
                        value={formData[resourceType][field.name] || ''}
                        onChange={(e) => handleInputChange(resourceType, field.name, e.target.value)}
                        required={field.required}
                      />
                    )}
                  </div>
                ))}
              </div>

              <div className="form-actions">
                <button
                  type="submit"
                  className="btn btn-primary"
                  disabled={loading[resourceType]}
                >
                  {loading[resourceType] ? '作成中...' : '作成'}
                </button>
                <button
                  type="button"
                  className="btn btn-secondary"
                  onClick={() => clearForm(resourceType)}
                >
                  クリア
                </button>
              </div>
            </form>
          </div>

          {/* データテーブル */}
          <div className="data-section">
            <h4 className="data-section-title">データ一覧</h4>
            {data.length > 0 ? (
              <table className="table">
                <thead>
                  <tr>
                    {Object.keys(data[0]).map((key) => (
                      <th key={key}>{key}</th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {data.map((item, index) => (
                    <tr key={index}>
                      {Object.values(item).map((value: any, idx) => (
                        <td key={idx}>{String(value)}</td>
                      ))}
                    </tr>
                  ))}
                </tbody>
              </table>
            ) : (
              <p className="text-muted text-center mt-3">データがありません</p>
            )}
          </div>
        </div>
      </div>
    );
  };

  return (
    <div className="container">
      <div className="header">
        <h1>Mist ED Backend デバッグツール</h1>
        <p>バックエンドの動作確認用データを管理できます</p>
      </div>

      <div className="d-flex gap-2 mb-3">
        <button className="btn btn-success" onClick={handleSeedData}>
          シードデータ作成
        </button>
        <button className="btn btn-danger" onClick={handleResetDatabase}>
          データベースリセット
        </button>
      </div>

      <div className="grid">
        {renderDataTable('組織', 'organizations', organizations)}
        {renderDataTable('ユーザー', 'users', users)}
        {renderDataTable('部屋', 'rooms', rooms)}
        {renderDataTable('デバイス', 'devices', devices)}
        {renderDataTable('科目', 'subjects', subjects)}
        {renderDataTable('授業', 'lessons', lessons)}
      </div>

      {notification.show && (
        <div className={`notification notification-${notification.type}`}>
          {notification.message}
        </div>
      )}
    </div>
  );
};

export default App;
