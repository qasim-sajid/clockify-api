CREATE TABLE IF NOT EXISTS public."user"
(
    _id character varying COLLATE pg_catalog."default" NOT NULL,
    name character varying COLLATE pg_catalog."default" NOT NULL,
    email character varying COLLATE pg_catalog."default" NOT NULL,
    username character varying COLLATE pg_catalog."default" NOT NULL,
    password character varying COLLATE pg_catalog."default",
    CONSTRAINT user_pkey PRIMARY KEY (_id),
    CONSTRAINT email_unique UNIQUE (email)
        INCLUDE(email),
    CONSTRAINT username_unique UNIQUE (username)
        INCLUDE(username)
);

CREATE TABLE IF NOT EXISTS public.client
(
    _id character varying COLLATE pg_catalog."default" NOT NULL,
    name character varying COLLATE pg_catalog."default" NOT NULL,
    address character varying COLLATE pg_catalog."default",
    note character varying COLLATE pg_catalog."default",
    is_archived boolean NOT NULL,
    CONSTRAINT client_pkey PRIMARY KEY (_id)
);

CREATE TABLE IF NOT EXISTS public.workspace
(
    _id character varying COLLATE pg_catalog."default" NOT NULL,
    name character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT workspace_pkey PRIMARY KEY (_id)
);

CREATE TABLE IF NOT EXISTS public.team_role
(
    _id character varying COLLATE pg_catalog."default" NOT NULL,
    role character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT team_role_pkey PRIMARY KEY (_id)
);

CREATE TABLE IF NOT EXISTS public.team_member
(
    _id character varying COLLATE pg_catalog."default" NOT NULL,
    billable_rate numeric NOT NULL,
    workspace_id character varying COLLATE pg_catalog."default" NOT NULL,
    user_email character varying COLLATE pg_catalog."default" NOT NULL,
    team_role_id character varying COLLATE pg_catalog."default",
    CONSTRAINT team_member_pkey PRIMARY KEY (_id),
    CONSTRAINT team_member_team_role_id_fkey FOREIGN KEY (team_role_id)
        REFERENCES public.team_role (_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT team_member_user_email_fkey FOREIGN KEY (user_email)
        REFERENCES public."user" (email) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT team_member_workspace_id_fkey FOREIGN KEY (workspace_id)
        REFERENCES public.workspace (_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

CREATE TABLE IF NOT EXISTS public.team_group
(
    _id character varying COLLATE pg_catalog."default" NOT NULL,
    name character varying COLLATE pg_catalog."default" NOT NULL,
    workspace_id character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT team_group_pkey PRIMARY KEY (_id),
    CONSTRAINT team_group_workspace_id_fkey FOREIGN KEY (workspace_id)
        REFERENCES public.workspace (_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

CREATE TABLE IF NOT EXISTS public.tag
(
    _id character varying COLLATE pg_catalog."default" NOT NULL,
    name character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT tag_pkey PRIMARY KEY (_id)
);

CREATE TABLE IF NOT EXISTS public.project
(
    _id character varying COLLATE pg_catalog."default" NOT NULL,
    name character varying COLLATE pg_catalog."default" NOT NULL,
    color_tag character varying COLLATE pg_catalog."default" NOT NULL,
    is_public boolean NOT NULL,
    tracked_hours numeric NOT NULL,
    tracked_amount numeric NOT NULL,
    progress_percentage numeric NOT NULL,
    client_id character varying COLLATE pg_catalog."default",
    workspace_id character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT project_pkey PRIMARY KEY (_id),
    CONSTRAINT project_client_id_fkey FOREIGN KEY (client_id)
        REFERENCES public.client (_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT project_workspace_id_fkey FOREIGN KEY (workspace_id)
        REFERENCES public.workspace (_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

CREATE TABLE IF NOT EXISTS public.task
(
    _id character varying COLLATE pg_catalog."default" NOT NULL,
    description character varying COLLATE pg_catalog."default",
    billable boolean NOT NULL,
    start_time character varying COLLATE pg_catalog."default" NOT NULL,
    end_time character varying COLLATE pg_catalog."default",
    date character varying COLLATE pg_catalog."default" NOT NULL,
    is_active boolean NOT NULL,
    project_id character varying COLLATE pg_catalog."default",
    CONSTRAINT task_pkey PRIMARY KEY (_id),
    CONSTRAINT task_project_id_fkey FOREIGN KEY (project_id)
        REFERENCES public.project (_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

CREATE TABLE IF NOT EXISTS public.project_team_group
(
    project_id character varying COLLATE pg_catalog."default" NOT NULL,
    team_group_id character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT project_team_group_project_id_fkey FOREIGN KEY (project_id)
        REFERENCES public.project (_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT project_team_group_team_group_id_fkey FOREIGN KEY (team_group_id)
        REFERENCES public.team_group (_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID
);

CREATE TABLE IF NOT EXISTS public.project_team_member
(
    project_id character varying COLLATE pg_catalog."default" NOT NULL,
    team_member_id character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT project_team_member_project_id_fkey FOREIGN KEY (project_id)
        REFERENCES public.project (_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT project_team_member_team_member_id_fkey FOREIGN KEY (team_member_id)
        REFERENCES public.team_member (_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID
);

CREATE TABLE IF NOT EXISTS public.team_group_team_member
(
    team_group_id character varying COLLATE pg_catalog."default" NOT NULL,
    team_member_id character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT team_group_id_fkey FOREIGN KEY (team_group_id)
        REFERENCES public.team_group (_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT team_member_id_fkey FOREIGN KEY (team_member_id)
        REFERENCES public.team_member (_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID
);

CREATE TABLE IF NOT EXISTS public.task_tag
(
    task_id character varying COLLATE pg_catalog."default" NOT NULL,
    tag_id character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT task_tag_tag_id_fkey FOREIGN KEY (tag_id)
        REFERENCES public.tag (_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT task_tag_task_id_fkey FOREIGN KEY (task_id)
        REFERENCES public.task (_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID
);