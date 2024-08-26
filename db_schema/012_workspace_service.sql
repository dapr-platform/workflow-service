-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE o_workflow(
                           id VARCHAR(255) NOT NULL,
                           created_by VARCHAR(32) NOT NULL,
                           created_time TIMESTAMP NOT NULL,
                           updated_by VARCHAR(32) NOT NULL,
                           updated_time TIMESTAMP NOT NULL,
                           name VARCHAR(255) NOT NULL,
                           status INTEGER NOT NULL,
                           content text NOT NULL,
                           description text NOT NULL,
                           type VARCHAR(255) NOT NULL,
                           workflow_id VARCHAR(255) NOT NULL,
                           running_id VARCHAR(255) NOT NULL,
                           PRIMARY KEY (id)
);

COMMENT ON TABLE o_workflow IS '流程表';
COMMENT ON COLUMN o_workflow.id IS '唯一标识';
COMMENT ON COLUMN o_workflow.created_time IS '创建时间';
COMMENT ON COLUMN o_workflow.updated_time IS '修改时间';
COMMENT ON COLUMN o_workflow.created_by IS '创建人id';
COMMENT ON COLUMN o_workflow.updated_by IS '修改人id';
COMMENT ON COLUMN o_workflow.name IS '名称';
COMMENT ON COLUMN o_workflow.status IS '状态 0:未发布1:已发布,2:已生效';
COMMENT ON COLUMN o_workflow.content IS '内容';
COMMENT ON COLUMN o_workflow.description IS '描述';
COMMENT ON COLUMN o_workflow.type IS '类型';
COMMENT ON COLUMN o_workflow.workflow_id IS 'temporal中id';
COMMENT ON COLUMN o_workflow.running_id IS '运行id';

CREATE TABLE o_workflow_trigger(
                                   id VARCHAR(255) NOT NULL,
                                   create_at TIMESTAMP NOT NULL,
                                   update_at TIMESTAMP NOT NULL,
                                   create_id VARCHAR(255) NOT NULL,
                                   update_id VARCHAR(255) NOT NULL,
                                   workflow_id VARCHAR(255) NOT NULL,
                                   type INTEGER NOT NULL,
                                   PRIMARY KEY (id)
);

COMMENT ON TABLE o_workflow_trigger IS '触发流程节点表';
COMMENT ON COLUMN o_workflow_trigger.id IS '唯一标识';
COMMENT ON COLUMN o_workflow_trigger.create_at IS '创建时间';
COMMENT ON COLUMN o_workflow_trigger.update_at IS '修改时间';
COMMENT ON COLUMN o_workflow_trigger.create_id IS '创建人id';
COMMENT ON COLUMN o_workflow_trigger.update_id IS '修改人id';
COMMENT ON COLUMN o_workflow_trigger.workflow_id IS '工作流id';
COMMENT ON COLUMN o_workflow_trigger.type IS '类型(1:systemEvent,2:cron,3:form,4:manual)';

CREATE TABLE o_workflow_trigger_system_event(
                                                pubsub varchar(255) not null,
                                                topic varchar(255) not null,
                                                filter_expr varchar(1024) not null
)INHERITS(o_workflow_trigger);

COMMENT ON TABLE o_workflow_trigger_system_event IS '触发流程节点表';
COMMENT ON COLUMN o_workflow_trigger_system_event.id IS '唯一标识';
COMMENT ON COLUMN o_workflow_trigger_system_event.create_at IS '创建时间';
COMMENT ON COLUMN o_workflow_trigger_system_event.update_at IS '修改时间';
COMMENT ON COLUMN o_workflow_trigger_system_event.create_id IS '创建人id';
COMMENT ON COLUMN o_workflow_trigger_system_event.update_id IS '修改人id';
COMMENT ON COLUMN o_workflow_trigger_system_event.workflow_id IS '工作流id';
COMMENT ON COLUMN o_workflow_trigger_system_event.type IS '类型(1:systemEvent,2:cron,3:form,4:manual)';
COMMENT ON COLUMN o_workflow_trigger_system_event.pubsub IS 'pubsub';
COMMENT ON COLUMN o_workflow_trigger_system_event.topic IS 'topic';
COMMENT ON COLUMN o_workflow_trigger_system_event.filter_expr IS '过滤表达式';

CREATE TABLE o_workflow_trigger_cron(
                                        cron varchar(255) not null

)INHERITS(o_workflow_trigger);

COMMENT ON TABLE o_workflow_trigger_cron IS '触发流程节点表';
COMMENT ON COLUMN o_workflow_trigger_cron.id IS '唯一标识';
COMMENT ON COLUMN o_workflow_trigger_cron.create_at IS '创建时间';
COMMENT ON COLUMN o_workflow_trigger_cron.update_at IS '修改时间';
COMMENT ON COLUMN o_workflow_trigger_cron.create_id IS '创建人id';
COMMENT ON COLUMN o_workflow_trigger_cron.update_id IS '修改人id';
COMMENT ON COLUMN o_workflow_trigger_cron.workflow_id IS '工作流id';
COMMENT ON COLUMN o_workflow_trigger_cron.type IS '类型(1:systemEvent,2:cron,3:form,4:manual)';
COMMENT ON COLUMN o_workflow_trigger_cron.cron IS 'cron表达式';


CREATE TABLE o_workflow_usertask_process(
                                            id VARCHAR(255) NOT NULL,
                                            create_at TIMESTAMP NOT NULL,
                                            update_at TIMESTAMP NOT NULL,
                                            create_id VARCHAR(255) NOT NULL,
                                            update_id VARCHAR(255) NOT NULL,
                                            run_workflow_id VARCHAR(255) NOT NULL,
                                            run_id VARCHAR(255) NOT NULL,
                                            to_type VARCHAR(255) NOT NULL,
                                            to_id VARCHAR(255) NOT NULL,
                                            notify_type VARCHAR(255) NOT NULL,
                                            form_ref_id VARCHAR(255) NOT NULL,
                                            form_def text NOT NULL,
                                            extra_form_def text NOT NULL,
                                            form_value text NOT NULL,
                                            token VARCHAR(1024) NOT NULL,
                                            node_id VARCHAR(255) NOT NULL,
                                            outgoing_properties text not null,
                                            PRIMARY KEY (id)
);

COMMENT ON TABLE o_workflow_usertask_process IS '工作流用户任务处理';
COMMENT ON COLUMN o_workflow_usertask_process.id IS '唯一标识';
COMMENT ON COLUMN o_workflow_usertask_process.create_at IS '创建时间';
COMMENT ON COLUMN o_workflow_usertask_process.update_at IS '修改时间';
COMMENT ON COLUMN o_workflow_usertask_process.create_id IS '创建人id';
COMMENT ON COLUMN o_workflow_usertask_process.update_id IS '修改人id';
COMMENT ON COLUMN o_workflow_usertask_process.run_workflow_id IS 'running workflow id';
COMMENT ON COLUMN o_workflow_usertask_process.run_id IS 'run id';
COMMENT ON COLUMN o_workflow_usertask_process.to_type IS 'to type(user,organization,role)';
COMMENT ON COLUMN o_workflow_usertask_process.to_id IS 'to id';
COMMENT ON COLUMN o_workflow_usertask_process.notify_type IS '通知类型(sms,email,event)';
COMMENT ON COLUMN o_workflow_usertask_process.form_ref_id IS 'form表单id';
COMMENT ON COLUMN o_workflow_usertask_process.form_def IS 'form表单附加的定义，例如只读';
COMMENT ON COLUMN o_workflow_usertask_process.extra_form_def IS '扩展form定义(json格式字符串)';
COMMENT ON COLUMN o_workflow_usertask_process.form_value IS '传递的值（json格式字符串）';
COMMENT ON COLUMN o_workflow_usertask_process.token IS 'workflow token []byte';
COMMENT ON COLUMN o_workflow_usertask_process.node_id IS 'json node id';
COMMENT ON COLUMN o_workflow_usertask_process.outgoing_properties IS 'user task的下级连接属性（json字符串)';




-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS p_workflow cascade;
DROP TABLE IF EXISTS p_workflow_trigger cascade ;
DROP TABLE IF EXISTS p_workflow_trigger_system_event;
DROP TABLE IF EXISTS p_workflow_trigger_cron;
DROP TABLE IF EXISTS p_workflow_usertask_process;

-- +goose StatementEnd
