CREATE DATABASE {{ .User }};
ALTER DATABASE {{ .User }} SET allow_snapshot_isolation ON;
ALTER DATABASE {{ .User }} SET SINGLE_USER WITH ROLLBACK IMMEDIATE;
ALTER DATABASE {{ .User }} SET read_committed_snapshot ON;
ALTER DATABASE {{ .User }} SET MULTI_USER;