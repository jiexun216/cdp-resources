package hook

import (
	"encoding/json"
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"strconv"
)

//1.# console
//monitoring-platform-access
//2.# cml
//associatedCRP
//4.# warehouse
//istio-injection
//5.# compute
//istio-injection
//6.# implala
//istio-injection
//7.# monitoring
//cdp.cloudera/version

//3.# cml-user
//ds-parent-namespace

//only the # cml-user//ds-parent-namespace//
// cdp-resources admission is to delete resource limit
func createSpecSecurityAndModifyResourcesContextPatch(pod corev1.Pod, availableAnnotations map[string]string, annotations map[string]string, availableLabels map[string]string, labels map[string]string) ([]byte, error) {
	var patch []patchOperation
	// update Annotation to set admissionWebhookAnnotationStatusKey: "mutated"
	patch = append(patch, updateAnnotation(availableAnnotations, annotations)...)
	// add labels
	patch = append(patch, updateLabels(availableLabels, labels)...)

	// delete pod resourse
	var size = len(pod.Spec.Containers)
	for i := 0; i < size; i++ {
		removeResourse := patchOperation{
			Op:    "replace",
			Path:  "/spec/containers/" + strconv.Itoa(i) + "/resources",
			Value: nil,
		}
		glog.Infof("remove  pod container resource for value: %v", removeResourse)
		patch = append(patch, removeResourse)
	}

	return json.Marshal(patch)
}




func updateAnnotation(target map[string]string, added map[string]string) (patch []patchOperation) {
	for key, value := range added {
		if target == nil || target[key] == "" {
			target = map[string]string{}
			patch = append(patch, patchOperation{
				Op:   "add",
				Path: "/metadata/annotations",
				Value: map[string]string{
					key: value,
				},
			})
		} else {
			patch = append(patch, patchOperation{
				Op:    "replace",
				Path:  "/metadata/annotations/" + key,
				Value: value,
			})
		}
	}
	return patch
}

func updateLabels(target map[string]string, added map[string]string) (patch []patchOperation) {
	values := make(map[string]string)
	for key, value := range added {
		if target == nil || target[key] == "" {
			values[key] = value
		}
	}
	patch = append(patch, patchOperation{
		Op:    "add",
		Path:  "/metadata/labels",
		Value: values,
	})
	return patch
}
