package Test

import (
	"github.com/mhthrh/GigaFileProcess/ftp"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFtpList(t *testing.T) {
	path := "myDir"
	client, err := ftp.New(ftpIp, ftpUser, ftpPassword, ftpPort)
	require.NoError(t, err)
	require.NotNil(t, client)

	err = client.CreateDirectory(path)
	require.NoError(t, err)

	err = client.DeleteDirectory(path)
	require.NoError(t, err)

	err = client.Upload(filetestName4Upload, filetestName4Upload, fileTestDirector)
	require.NoError(t, err)

	err = client.Download(filetestName4Upload, filetestName4Download, fileTestDirector)
	require.NoError(t, err)
	// compare hash of files, if same PASS
	err = client.Close()
	require.NoError(t, err)
}
