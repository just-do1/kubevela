import React, { useCallback, useEffect, useState } from 'react';

import { Button, Card, message, Popconfirm, Space, Table } from 'antd';
import { history, Link, useModel, useRequest } from 'umi';

import { deleteApplication, getApplications } from '@/services/application';
import { PlusOutlined } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-layout';

export default () => {
  const { currentEnvironment } = useModel('useEnvironmentModel');
  const [removingApps, setRemovingApps] = useState<string[]>([]);
  const { data, loading, run: loadApps } = useRequest(
    async () => {
      if (currentEnvironment?.envName == null) {
        return { code: 200, data: [] } as API.VelaResponse<Array<API.Application>>;
      }
      return getApplications(currentEnvironment.envName);
    },
    {
      refreshDeps: [currentEnvironment],
    },
  );

  // Make sure that the deleted app disappears before stopping polling
  useEffect(() => {
    // Return when removingApps is empty
    if (removingApps.length <= 0) {
      return undefined;
    }
    // If there is no data in the app list, clear it directly
    if (data == null || data.length === 0) {
      setRemovingApps([]);
      return undefined;
    }

    // polling load apps data
    const timer = setInterval(() => {
      loadApps();
    }, 500);

    // If the deleted data disappears, reset the value
    const appNames = data.map((i) => i.name);
    const newRemoveingApps = removingApps.filter((ra) => appNames.includes(ra));
    if (newRemoveingApps.join(',') !== removingApps.join(',')) {
      setRemovingApps(newRemoveingApps);
    }
    return () => {
      clearInterval(timer);
    };
  }, [data, removingApps]);

  // delete application
  const remove = useCallback(
    async (appName: string) => {
      const envName = currentEnvironment?.envName;
      if (envName == null) {
        throw new Error('Unable to determine the current environment name.');
      }

      // append removing app
      setRemovingApps(removingApps.concat([appName]));

      return deleteApplication(envName, appName);
    },
    [currentEnvironment],
  );

  return (
    <PageContainer>
      <Card>
        <div style={{ marginBottom: '10px' }}>
          <Space>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              href="/applications/create"
              onClick={(e) => {
                history.push('/applications/create');
                e.preventDefault();
              }}
            >
              Create
            </Button>
          </Space>
        </div>
        <Table
          dataSource={data ?? []}
          rowKey={(record) => record.name}
          loading={loading ? { delay: 300 } : undefined}
          columns={[
            {
              title: 'Name',
              dataIndex: 'name',
              key: 'name',
              render: (text, record) => {
                return (
                  <Link
                    to={{
                      pathname: `${window.routerBase}applications/${record.name}`,
                    }}
                  >
                    {text}
                  </Link>
                );
              },
            },
            {
              title: 'Status',
              dataIndex: 'status',
              key: 'status',
              render: (text) => {
                return text;
              },
            },
            {
              title: 'Created Time',
              dataIndex: 'createdTime',
              key: 'createdTime',
              render: (text) => {
                return text;
              },
            },
            {
              title: 'Actions',
              dataIndex: 'Actions',
              key: 'Actions',
              render: (text, { name }) => {
                return (
                  <Space>
                    <Popconfirm
                      title="Are you sure to delete this application?"
                      onConfirm={() => {
                        remove(name).then(({ code, data: content }) => {
                          if (code === 200) {
                            message.success({
                              content,
                              key: 'remove',
                            });
                          } else {
                            message.error({
                              content,
                              key: 'remove',
                            });
                          }
                        });
                      }}
                    >
                      <Button type="link" size="small" danger>
                        Delete
                      </Button>
                    </Popconfirm>
                  </Space>
                );
              },
            },
          ]}
        />
      </Card>
    </PageContainer>
  );
};
