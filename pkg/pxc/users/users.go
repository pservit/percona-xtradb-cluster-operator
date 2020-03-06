package users

import (
	api "github.com/percona/percona-xtradb-cluster-operator/pkg/apis/pxc/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Job(cr *api.PerconaXtraDBCluster) *batchv1.Job {
	/*labels := make(map[string]string)
	labels["type"] = "usermanager"
	labels["cluster"] = cr.Name*/
	return &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: cr.Namespace,
			//Labels:    labels,
		},
	}
}

func JobSpec(rootPass string, conns string, job *batchv1.Job) batchv1.JobSpec {
	backbackoffLimit := int32(1)
	return batchv1.JobSpec{
		BackoffLimit: &backbackoffLimit,
		Selector:     &metav1.LabelSelector{
			//MatchLabels: job.Labels,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				//Labels: job.Labels,
			},
			Spec: corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyNever,
				Containers: []corev1.Container{
					{
						Name:            "pxcusers",
						Image:           "nonemax/users:latest",
						ImagePullPolicy: corev1.PullAlways,
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/go/src/github.com/percona/pxcusers",
								Name:      "userssecret",
								ReadOnly:  true,
							},
						},
						Env: []corev1.EnvVar{
							{
								Name:  "PXC_CONNS",
								Value: conns,
							},
							{
								Name:  "PXC-ROOT-PASS",
								Value: rootPass,
							},
						},
					},
				},
				Volumes: []corev1.Volume{
					{
						Name: "userssecret",
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: "secret-for-users",
							},
						},
					},
				},
			},
		},
	}
}
