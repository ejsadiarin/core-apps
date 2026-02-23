"use client";

import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useCategories, useTags, usePriorityGroups } from "@/hooks/use-budget";
import type { CreateExpenseRequest, UpdateExpenseRequest } from "@/types/api";
import { Badge } from "@/components/ui/badge";
import { X } from "lucide-react";

interface ExpenseFormProps {
  initialData?: UpdateExpenseRequest & { id?: string };
  onSubmit: (data: CreateExpenseRequest | UpdateExpenseRequest) => void | Promise<void>;
  onCancel?: () => void;
}

export function ExpenseForm({ initialData, onSubmit, onCancel }: ExpenseFormProps) {
  const { data: categories } = useCategories();
  const { data: tags } = useTags();
  const { data: priorityGroups } = usePriorityGroups();

  const [formData, setFormData] = useState<CreateExpenseRequest>({
    description: initialData?.description || "",
    amount: initialData?.amount || 0,
    currency: initialData?.currency || "PHP",
    category_id: initialData?.category_id,
    priority_group_id: initialData?.priority_group_id,
    expense_date: initialData?.expense_date || new Date().toISOString().split('T')[0],
    notes: initialData?.notes,
    tag_ids: initialData?.tag_ids || [],
    recurring_type: initialData?.recurring_type || null,
    start_date: initialData?.start_date || undefined,
    end_date: initialData?.end_date || undefined,
  });

  const [selectedTags, setSelectedTags] = useState<string[]>(initialData?.tag_ids || []);
  const [noEndDate, setNoEndDate] = useState(!initialData?.end_date);

  useEffect(() => {
    if (initialData?.tag_ids) {
      setSelectedTags(initialData.tag_ids);
    }
  }, [initialData?.tag_ids]);

  const handleTagToggle = (tagId: string) => {
    setSelectedTags((prev) =>
      prev.includes(tagId)
        ? prev.filter((id) => id !== tagId)
        : [...prev, tagId]
    );
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await onSubmit({ ...formData, tag_ids: selectedTags });
  };

  const isRecurring = formData.recurring_type !== null && formData.recurring_type !== undefined;

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Description */}
      <div className="space-y-2">
        <Label htmlFor="description">Description *</Label>
        <Input
          id="description"
          required
          value={formData.description}
          onChange={(e) => setFormData({ ...formData, description: e.target.value })}
          placeholder="e.g., Groceries at Whole Foods"
        />
      </div>

      {/* Amount and Currency */}
      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="amount">Amount *</Label>
          <Input
            id="amount"
            type="number"
            step="0.01"
            required
            value={formData.amount || ""}
            onChange={(e) => setFormData({ ...formData, amount: parseFloat(e.target.value) || 0 })}
            placeholder="0.00"
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="currency">Currency</Label>
          <Input
            id="currency"
            maxLength={3}
            value={formData.currency}
            onChange={(e) => setFormData({ ...formData, currency: e.target.value.toUpperCase() })}
            placeholder="PHP"
          />
        </div>
      </div>

      {/* Category */}
      <div className="space-y-2">
        <Label htmlFor="category">Category</Label>
        <Select
          value={formData.category_id || "none"}
          onValueChange={(value) =>
            setFormData({ ...formData, category_id: value === "none" ? undefined : value })
          }
        >
          <SelectTrigger id="category">
            <SelectValue placeholder="Select a category" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="none">None</SelectItem>
            {categories?.map((category) => (
              <SelectItem key={category.id} value={category.id}>
                {category.icon && <span className="mr-2">{category.icon}</span>}
                {category.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {/* Priority Group (Need/Want/Savings) */}
      <div className="space-y-2">
        <Label htmlFor="priority_group">Priority (50/30/20)</Label>
        <Select
          value={formData.priority_group_id || "none"}
          onValueChange={(value) =>
            setFormData({ ...formData, priority_group_id: value === "none" ? undefined : value })
          }
        >
          <SelectTrigger id="priority_group">
            <SelectValue placeholder="Select priority" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="none">Unclassified</SelectItem>
            {priorityGroups?.map((pg) => (
              <SelectItem key={pg.id} value={pg.id}>
                {pg.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {/* Recurring Type */}
      <div className="space-y-2">
        <Label htmlFor="recurring_type">Type</Label>
        <Select
          value={formData.recurring_type || "one-time"}
          onValueChange={(value) =>
            setFormData({
              ...formData,
              recurring_type: value === "one-time" ? null : (value as "daily" | "weekly" | "monthly" | "yearly"),
              start_date: value !== "one-time" ? formData.expense_date : undefined,
              end_date: value !== "one-time" ? formData.end_date : undefined,
            })
          }
        >
          <SelectTrigger id="recurring_type">
            <SelectValue placeholder="Select type" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="one-time">One-time</SelectItem>
            <SelectItem value="daily">Daily Recurring</SelectItem>
            <SelectItem value="weekly">Weekly Recurring</SelectItem>
            <SelectItem value="monthly">Monthly Recurring</SelectItem>
            <SelectItem value="yearly">Yearly Recurring</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Date for one-time, Start Date for recurring */}
      <div className="space-y-2">
        <Label htmlFor="expense_date">{isRecurring ? "Start Date *" : "Date *"}</Label>
        <Input
          id="expense_date"
          type="date"
          required
          value={formData.expense_date}
          onChange={(e) => {
            const newDate = e.target.value;
            setFormData({
              ...formData,
              expense_date: newDate,
              start_date: isRecurring ? newDate : formData.start_date,
            });
          }}
        />
      </div>

      {/* End Date for recurring expenses */}
      {isRecurring && (
        <>
          <div className="flex items-center space-x-2">
            <input
              type="checkbox"
              id="no_end_date"
              checked={noEndDate}
              onChange={(e) => {
                setNoEndDate(e.target.checked);
                if (e.target.checked) {
                  setFormData({ ...formData, end_date: undefined });
                }
              }}
              className="h-4 w-4 rounded border-gray-300"
            />
            <Label htmlFor="no_end_date" className="font-normal cursor-pointer">
              No end date (ongoing)
            </Label>
          </div>

          {!noEndDate && (
            <div className="space-y-2">
              <Label htmlFor="end_date">End Date</Label>
              <Input
                id="end_date"
                type="date"
                value={formData.end_date || ""}
                min={formData.start_date || formData.expense_date}
                onChange={(e) => setFormData({ ...formData, end_date: e.target.value || undefined })}
              />
              {formData.end_date && formData.start_date && formData.end_date < formData.start_date && (
                <p className="text-sm text-red-600">End date must be on or after start date</p>
              )}
            </div>
          )}

          {initialData?.end_date === undefined && noEndDate && (
            <div className="text-sm text-muted-foreground bg-blue-50 border border-blue-200 rounded px-3 py-2">
              <span className="font-medium">Ongoing</span> - This expense will continue indefinitely
            </div>
          )}
        </>
      )}

      {/* Tags */}
      <div className="space-y-2">
        <Label>Tags</Label>
        <div className="flex flex-wrap gap-2 p-2 border rounded-md min-h-[60px]">
          {tags?.map((tag) => {
            const isSelected = selectedTags.includes(tag.id);
            return (
              <Badge
                key={tag.id}
                variant={isSelected ? "default" : "outline"}
                className="cursor-pointer"
                style={{
                  backgroundColor: isSelected && tag.color ? tag.color : undefined,
                  borderColor: tag.color || undefined,
                }}
                onClick={() => handleTagToggle(tag.id)}
              >
                {tag.name}
                {isSelected && <X className="ml-1 h-3 w-3" />}
              </Badge>
            );
          })}
        </div>
      </div>

      {/* Notes */}
      <div className="space-y-2">
        <Label htmlFor="notes">Notes</Label>
        <Textarea
          id="notes"
          value={formData.notes || ""}
          onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
          placeholder="Additional notes about this expense..."
          rows={3}
        />
      </div>

      {/* Actions */}
      <div className="flex gap-3 justify-end">
        {onCancel && (
          <Button type="button" variant="outline" onClick={onCancel}>
            Cancel
          </Button>
        )}
        <Button type="submit">
          {initialData?.id ? "Update" : "Create"} Expense
        </Button>
      </div>
    </form>
  );
}
