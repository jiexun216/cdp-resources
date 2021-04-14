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
func createDeleteResourceContextPatch(pod corev1.Pod, availableAnnotations map[string]string, annotations map[string]string) ([]byte, error) {
	var patch []patchOperation
	// update Annotation to set admissionWebhookAnnotationStatusKey: "mutated"
	patch = append(patch, updateAnnotation(availableAnnotations, annotations)...)

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

